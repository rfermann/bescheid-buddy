package paper_letters_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/bored-bots/bored-template/internal/repository/cases_repository"
	"github.com/bored-bots/bored-template/internal/services/reply_service"
	"github.com/bored-bots/bored-template/internal/services/users_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrGenerationNotFound      = errors.New("paper letter generation was not found")
	ErrInvalidInput            = errors.New("paper letter generation input is invalid")
	ErrUnsupportedMode         = errors.New("unsupported paper letter generation mode")
	ErrInvalidSource           = errors.New("paper letter source is invalid")
	ErrSourceNotAnalyzed       = errors.New("paper letter source analysis is not completed")
	ErrUnsupportedModel        = errors.New("unsupported paper letter model")
	ErrUnsupportedLocale       = errors.New("unsupported paper letter locale")
	ErrInsufficientBuddyPoints = users_service.ErrInsufficientBuddyPoints
)

const generationRequestConstraint = "document_reply_generations_user_client_request_unique"

type GenerationQueue interface {
	QueueDocumentReplyGenerationJobTx(ctx context.Context, tx pgx.Tx, replyGenerationID uuid.UUID) (int64, error)
	QueuePaperLetterDeliveryJobTx(ctx context.Context, tx pgx.Tx, deliveryID uuid.UUID) (int64, error)
}

type ModelCatalog interface {
	GetModelByProviderModelID(ctx context.Context, providerModelID string) (*db.AiModel, error)
}

type PointSpender interface {
	SpendTx(ctx context.Context, tx pgx.Tx, input users_service.SpendTxInput) (*users_service.SpendTxResult, error)
}

type Service struct {
	repository  cases_repository.Repository
	queue       GenerationQueue
	models      ModelCatalog
	points      PointSpender
	storage     PaperLetterStorage
	maxPDFBytes int64
	pingen      config.PingenConfig
}

func New(repository cases_repository.Repository, queue GenerationQueue, models ModelCatalog, points PointSpender, storage PaperLetterStorage, maxPDFBytes int64, pingen config.PingenConfig) *Service {
	return &Service{repository: repository, queue: queue, models: models, points: points, storage: storage, maxPDFBytes: maxPDFBytes, pingen: pingen}
}

type GenerateInput struct {
	UserID           uuid.UUID
	Mode             db.DocumentReplyGenerationMode
	SourceCaseID     *uuid.UUID
	SourceDocumentID *uuid.UUID
	ClientRequestID  *uuid.UUID
	Model            string
	Locale           string
	Input            reply_service.GenerationInput
}

type GenerateResult struct {
	*GenerationResult
	Existing bool
}

type GenerationResult struct {
	Generation       *db.DocumentReplyGeneration
	Input            reply_service.GenerationInput
	SourceCaseID     *uuid.UUID
	SourceDocumentID *uuid.UUID
}

func (s *Service) Generate(ctx context.Context, input GenerateInput) (*GenerateResult, error) {
	clientRequestID := uuid.New()
	if input.ClientRequestID != nil {
		if *input.ClientRequestID == uuid.Nil {
			return nil, fmt.Errorf("%w: client request ID", ErrInvalidInput)
		}
		clientRequestID = *input.ClientRequestID
		if existing, err := s.generationByRequest(ctx, input.UserID, clientRequestID); err == nil {
			decoded, err := s.decodeGeneration(ctx, existing)
			if err != nil {
				return nil, err
			}
			return &GenerateResult{GenerationResult: decoded, Existing: true}, nil
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}
	mode := input.Mode
	if mode == "" {
		mode = db.DocumentReplyGenerationModeGenerate
	}
	if mode != db.DocumentReplyGenerationModeGenerate && mode != db.DocumentReplyGenerationModeImprove {
		return nil, ErrUnsupportedMode
	}

	modelID := strings.TrimSpace(input.Model)
	model, err := s.models.GetModelByProviderModelID(ctx, modelID)
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && !model.Enabled) {
		return nil, ErrUnsupportedModel
	}
	if err != nil {
		return nil, fmt.Errorf("load generation model: %w", err)
	}

	locale := strings.ToLower(strings.TrimSpace(input.Locale))
	if !supportedLocale(locale) {
		return nil, ErrUnsupportedLocale
	}

	normalized := normalizeInput(input.Input)
	normalized.Locale = locale
	if err := validateInput(mode, normalized, input.SourceDocumentID != nil); err != nil {
		return nil, err
	}
	if err := s.validateSource(ctx, input.UserID, input.SourceCaseID, input.SourceDocumentID); err != nil {
		return nil, err
	}

	encodedInput, err := json.Marshal(normalized)
	if err != nil {
		return nil, fmt.Errorf("encode generation input: %w", err)
	}

	generationID := uuid.New()
	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin generation transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := s.points.SpendTx(ctx, tx, users_service.SpendTxInput{
		UserID:            input.UserID,
		Amount:            reply_service.GenerationCost,
		Description:       "Papierbrief generation",
		RelatedEntityID:   generationID,
		RelatedEntityType: db.RelatedEntityTypePaperLetterGeneration,
	}); err != nil {
		return nil, fmt.Errorf("spend generation BuddyPoints: %w", err)
	}

	txRepo := s.repository.WithTx(tx)
	documentID := pgtype.UUID{}
	if input.SourceDocumentID != nil {
		documentID = pgtype.UUID{Bytes: *input.SourceDocumentID, Valid: true}
	}
	generation, err := txRepo.CreateDocumentReplyGeneration(ctx, db.CreateDocumentReplyGenerationParams{
		ID:              generationID,
		UserID:          input.UserID,
		DocumentID:      documentID,
		Mode:            mode,
		ClientRequestID: clientRequestID,
		JobID:           pgtype.Int8{},
		Model:           model.ProviderModelID,
		Provider:        model.Provider,
		Language:        locale,
		Input:           encodedInput,
	})
	if err != nil {
		if isRequestConflict(err) {
			_ = tx.Rollback(ctx)
			existing, loadErr := s.generationByRequest(ctx, input.UserID, clientRequestID)
			if loadErr != nil {
				return nil, fmt.Errorf("load concurrent generation request: %w", loadErr)
			}
			decoded, decodeErr := s.decodeGeneration(ctx, existing)
			if decodeErr != nil {
				return nil, decodeErr
			}
			return &GenerateResult{GenerationResult: decoded, Existing: true}, nil
		}
		return nil, fmt.Errorf("create generation: %w", err)
	}

	jobID, err := s.queue.QueueDocumentReplyGenerationJobTx(ctx, tx, generationID)
	if err != nil {
		return nil, fmt.Errorf("enqueue generation: %w", err)
	}
	if err := txRepo.SetDocumentReplyGenerationJobID(ctx, db.SetDocumentReplyGenerationJobIDParams{
		ID:    generationID,
		JobID: pgtype.Int8{Int64: jobID, Valid: true},
	}); err != nil {
		return nil, fmt.Errorf("store generation job ID: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit generation: %w", err)
	}
	generation.JobID = pgtype.Int8{Int64: jobID, Valid: true}
	return &GenerateResult{GenerationResult: &GenerationResult{
		Generation:       generation,
		Input:            normalized,
		SourceCaseID:     copyUUID(input.SourceCaseID),
		SourceDocumentID: copyUUID(input.SourceDocumentID),
	}}, nil
}

func (s *Service) GetGeneration(ctx context.Context, userID, generationID uuid.UUID) (*GenerationResult, error) {
	generation, err := s.repository.GetDocumentReplyGenerationByIDAndUserID(ctx, db.GetDocumentReplyGenerationByIDAndUserIDParams{
		ID: generationID, UserID: userID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrGenerationNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("load generation: %w", err)
	}
	return s.decodeGeneration(ctx, generation)
}

func (s *Service) GetSourceGeneration(ctx context.Context, userID, caseID, documentID, generationID uuid.UUID) (*GenerationResult, error) {
	if _, err := s.repository.GetCaseByIDAndUserID(ctx, db.GetCaseByIDAndUserIDParams{ID: caseID, UserID: userID}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrGenerationNotFound
		}
		return nil, fmt.Errorf("load generation source case: %w", err)
	}
	result, err := s.GetGeneration(ctx, userID, generationID)
	if err != nil {
		return nil, err
	}
	if result.SourceCaseID == nil || *result.SourceCaseID != caseID ||
		result.SourceDocumentID == nil || *result.SourceDocumentID != documentID {
		return nil, ErrGenerationNotFound
	}
	return result, nil
}

func (s *Service) generationByRequest(ctx context.Context, userID, clientRequestID uuid.UUID) (*db.DocumentReplyGeneration, error) {
	return s.repository.GetDocumentReplyGenerationByUserAndClientRequestID(ctx, db.GetDocumentReplyGenerationByUserAndClientRequestIDParams{
		UserID: userID, ClientRequestID: clientRequestID,
	})
}

func (s *Service) validateSource(ctx context.Context, userID uuid.UUID, caseID, documentID *uuid.UUID) error {
	if (caseID == nil) != (documentID == nil) {
		return ErrInvalidSource
	}
	if caseID == nil {
		return nil
	}
	if _, err := s.repository.GetCaseByIDAndUserID(ctx, db.GetCaseByIDAndUserIDParams{ID: *caseID, UserID: userID}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInvalidSource
		}
		return fmt.Errorf("load source case: %w", err)
	}
	document, err := s.repository.GetCaseDocumentByID(ctx, db.GetCaseDocumentByIDParams{ID: *documentID})
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && document.CaseID != *caseID) {
		return ErrInvalidSource
	}
	if err != nil {
		return fmt.Errorf("load source document: %w", err)
	}
	analysis, err := s.repository.GetDocumentAnalysisByDocumentID(ctx, db.GetDocumentAnalysisByDocumentIDParams{DocumentID: *documentID})
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && analysis.Status != db.DocumentAnalysisStatusCompleted) {
		return ErrSourceNotAnalyzed
	}
	if err != nil {
		return fmt.Errorf("load source analysis: %w", err)
	}
	return nil
}

func (s *Service) decodeGeneration(ctx context.Context, generation *db.DocumentReplyGeneration) (*GenerationResult, error) {
	var input reply_service.GenerationInput
	if err := json.Unmarshal(generation.Input, &input); err != nil {
		return nil, fmt.Errorf("decode generation input: %w", err)
	}
	result := &GenerationResult{Generation: generation, Input: input}
	if generation.DocumentID.Valid {
		documentID := uuid.UUID(generation.DocumentID.Bytes)
		document, err := s.repository.GetCaseDocumentByID(ctx, db.GetCaseDocumentByIDParams{ID: documentID})
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("load generation source: %w", err)
		}
		if err == nil {
			caseID := document.CaseID
			result.SourceCaseID = &caseID
			result.SourceDocumentID = &documentID
		}
	}
	return result, nil
}

func copyUUID(value *uuid.UUID) *uuid.UUID {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func normalizeInput(input reply_service.GenerationInput) reply_service.GenerationInput {
	input.Sender = normalizeParty(input.Sender)
	input.Recipient = normalizeParty(input.Recipient)
	input.Category = strings.TrimSpace(input.Category)
	input.Subcategory = strings.TrimSpace(input.Subcategory)
	input.Subject = strings.TrimSpace(input.Subject)
	input.Reference = strings.TrimSpace(input.Reference)
	input.Instruction = strings.TrimSpace(input.Instruction)
	input.CurrentGermanLetter = strings.TrimSpace(input.CurrentGermanLetter)
	input.Feedback = strings.TrimSpace(input.Feedback)
	return input
}

func normalizeParty(party *reply_service.PartyInput) *reply_service.PartyInput {
	if party == nil {
		return nil
	}
	return &reply_service.PartyInput{Name: strings.TrimSpace(party.Name), Address: strings.TrimSpace(party.Address)}
}

func validateInput(mode db.DocumentReplyGenerationMode, input reply_service.GenerationInput, hasSource bool) error {
	fields := []struct {
		name  string
		value string
		max   int
	}{
		{"sender.name", partyName(input.Sender), 200},
		{"sender.address", partyAddress(input.Sender), 1000},
		{"recipient.name", partyName(input.Recipient), 200},
		{"recipient.address", partyAddress(input.Recipient), 1000},
		{"category", input.Category, 100},
		{"subcategory", input.Subcategory, 100},
		{"subject", input.Subject, 500},
		{"reference", input.Reference, 200},
		{"instruction", input.Instruction, 2500},
		{"currentGermanLetter", input.CurrentGermanLetter, 24000},
		{"feedback", input.Feedback, 2500},
	}
	for _, field := range fields {
		if utf8.RuneCountInString(field.value) > field.max {
			return fmt.Errorf("%w: %s exceeds %d characters", ErrInvalidInput, field.name, field.max)
		}
	}
	if input.Category != "" && !validCategory(input.Category) {
		return fmt.Errorf("%w: category", ErrInvalidInput)
	}
	if input.Subcategory != "" && !validSubcategory(input.Subcategory) {
		return fmt.Errorf("%w: subcategory", ErrInvalidInput)
	}
	if mode == db.DocumentReplyGenerationModeImprove {
		if input.CurrentGermanLetter == "" || input.Feedback == "" {
			return fmt.Errorf("%w: improve requires current German letter and feedback", ErrInvalidInput)
		}
		return nil
	}
	if !hasSource {
		if !completeParty(input.Sender) || !completeParty(input.Recipient) {
			return fmt.Errorf("%w: manual sender and recipient require name and address", ErrInvalidInput)
		}
		if input.Instruction == "" && input.Subject == "" && input.Reference == "" && input.Subcategory == "" {
			return fmt.Errorf("%w: manual generation needs letter context", ErrInvalidInput)
		}
	}
	return nil
}

func supportedLocale(locale string) bool {
	switch locale {
	case "en", "de", "ru", "ua", "pl", "ar":
		return true
	default:
		return false
	}
}

func validCategory(category string) bool {
	switch category {
	case "Jobcenter", "Krankenkasse", "Finanzamt", "Ausländerbehörde", "Rentenversicherung", "Arbeitsagentur",
		"Inkasso", "Gericht", "Schule/Kita", "Polizei", "Versicherung", "Bank", "Telekom/Internet", "Energie",
		"Other", "Universal":
		return true
	default:
		return false
	}
}

func validSubcategory(subcategory string) bool {
	switch subcategory {
	case "widerspruch", "fristverlaengerung", "unterlagen_nachreichen", "zahlung_klaeren", "termin_verschieben",
		"kuendigung", "freie_formulierung":
		return true
	default:
		return false
	}
}

func partyName(party *reply_service.PartyInput) string {
	if party == nil {
		return ""
	}
	return party.Name
}

func partyAddress(party *reply_service.PartyInput) string {
	if party == nil {
		return ""
	}
	return party.Address
}

func completeParty(party *reply_service.PartyInput) bool {
	return party != nil && party.Name != "" && party.Address != ""
}

func isRequestConflict(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == generationRequestConstraint
}
