package reply_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/bored-bots/bored-template/internal/services/llm_service"
	"github.com/bored-bots/bored-template/internal/services/users_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const GenerationCost int32 = 10

var (
	ErrGenerationNotFound = errors.New("reply generation was not found")
	ErrInvalidInput       = errors.New("reply generation input is invalid")
	ErrExtractionNotReady = errors.New("document extraction is not completed")
	ErrAnalysisNotReady   = errors.New("document analysis is not completed")
	ErrInvalidOutput      = errors.New("reply model did not return valid output")
)

type Repository interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
	GetDocumentReplyGenerationByID(ctx context.Context, params db.GetDocumentReplyGenerationByIDParams) (*db.DocumentReplyGeneration, error)
	MarkDocumentReplyGenerationProcessing(ctx context.Context, params db.MarkDocumentReplyGenerationProcessingParams) (*db.DocumentReplyGeneration, error)
	CompleteDocumentReplyGeneration(ctx context.Context, params db.CompleteDocumentReplyGenerationParams) (*db.DocumentReplyGeneration, error)
	FailActiveDocumentReplyGenerationTx(ctx context.Context, tx pgx.Tx, params db.FailActiveDocumentReplyGenerationParams) (*db.FailActiveDocumentReplyGenerationRow, error)
	GetDocumentExtractionTextByDocumentID(ctx context.Context, params db.GetDocumentExtractionTextByDocumentIDParams) (*db.GetDocumentExtractionTextByDocumentIDRow, error)
	GetDocumentAnalysisByDocumentID(ctx context.Context, params db.GetDocumentAnalysisByDocumentIDParams) (*db.DocumentAnalysisResult, error)
}

type PointRefunder interface {
	RefundTx(ctx context.Context, tx pgx.Tx, input users_service.RefundTxInput) (*users_service.RefundTxResult, error)
}

type LLM interface {
	Complete(ctx context.Context, model llm_service.Model, system, user string) (string, error)
}

type ModelCatalog interface {
	GetModelByProviderModelID(ctx context.Context, providerModelID string) (*db.AiModel, error)
}

type Service struct {
	repository Repository
	models     ModelCatalog
	llm        LLM
	points     PointRefunder
}

type Result struct {
	GermanLetter       string `json:"germanLetter"`
	UserLanguageLetter string `json:"userLanguageLetter"`
	Explanation        string `json:"explanation"`
	Subject            string `json:"subject"`
	Tone               string `json:"tone"`
}

type PartyInput struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}

type GenerationInput struct {
	Sender              *PartyInput `json:"sender,omitempty"`
	Recipient           *PartyInput `json:"recipient,omitempty"`
	Category            string      `json:"category,omitempty"`
	Subcategory         string      `json:"subcategory,omitempty"`
	Subject             string      `json:"subject,omitempty"`
	Reference           string      `json:"reference,omitempty"`
	Instruction         string      `json:"instruction"`
	CurrentGermanLetter string      `json:"currentGermanLetter,omitempty"`
	Feedback            string      `json:"feedback,omitempty"`
	Locale              string      `json:"locale,omitempty"`
}

func New(repository Repository, models ModelCatalog, llm LLM, points PointRefunder) *Service {
	return &Service{repository: repository, models: models, llm: llm, points: points}
}

func (s *Service) Generate(ctx context.Context, replyGenerationID uuid.UUID) error {
	generation, err := s.repository.GetDocumentReplyGenerationByID(ctx, db.GetDocumentReplyGenerationByIDParams{ID: replyGenerationID})
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrGenerationNotFound
	}
	if err != nil {
		return err
	}
	if generation.Status == db.DocumentReplyGenerationStatusCompleted || generation.Status == db.DocumentReplyGenerationStatusFailed {
		return nil
	}

	var generationInput GenerationInput
	if err := json.Unmarshal(generation.Input, &generationInput); err != nil {
		return s.FailAndRefund(ctx, replyGenerationID, ErrInvalidInput)
	}
	catalogModel, err := s.models.GetModelByProviderModelID(ctx, generation.Model)
	if errors.Is(err, pgx.ErrNoRows) {
		return s.FailAndRefund(ctx, replyGenerationID, llm_service.ErrUnsupportedModel)
	}
	if err != nil {
		return err
	}
	model := llm_service.Model{ID: catalogModel.ProviderModelID, Provider: catalogModel.Provider}
	if _, err := s.repository.MarkDocumentReplyGenerationProcessing(ctx, db.MarkDocumentReplyGenerationProcessingParams{ID: replyGenerationID}); err != nil {
		return err
	}

	var sourceText string
	var analysisMetadata []byte
	if generation.DocumentID.Valid {
		documentID := uuid.UUID(generation.DocumentID.Bytes)
		extraction, err := s.repository.GetDocumentExtractionTextByDocumentID(ctx, db.GetDocumentExtractionTextByDocumentIDParams{DocumentID: documentID})
		if errors.Is(err, pgx.ErrNoRows) || (err == nil && extraction.Status != db.DocumentAnalysisStatusCompleted) {
			return s.FailAndRefund(ctx, replyGenerationID, ErrExtractionNotReady)
		}
		if err != nil {
			return err
		}

		analysis, err := s.repository.GetDocumentAnalysisByDocumentID(ctx, db.GetDocumentAnalysisByDocumentIDParams{DocumentID: documentID})
		if errors.Is(err, pgx.ErrNoRows) || (err == nil && analysis.Status != db.DocumentAnalysisStatusCompleted) {
			return s.FailAndRefund(ctx, replyGenerationID, ErrAnalysisNotReady)
		}
		if err != nil {
			return err
		}
		sourceText = extraction.Text.String
		analysisMetadata = analysis.AnalysisMetadata
	} else if generation.Mode == db.DocumentReplyGenerationModeGenerate && !hasSubmittedContext(generationInput) {
		return s.FailAndRefund(ctx, replyGenerationID, ErrInvalidInput)
	}
	if generation.Mode == db.DocumentReplyGenerationModeImprove &&
		(strings.TrimSpace(generationInput.CurrentGermanLetter) == "" || strings.TrimSpace(generationInput.Feedback) == "") {
		return s.FailAndRefund(ctx, replyGenerationID, ErrInvalidInput)
	}

	system, user := buildPrompt(
		generation.Mode,
		generation.Language,
		generationInput,
		sourceText,
		analysisMetadata,
	)
	raw, err := s.llm.Complete(ctx, model, system, user)
	if err != nil {
		if IsPermanent(err) {
			return s.FailAndRefund(ctx, replyGenerationID, err)
		}
		return err
	}

	result, ok := decodeResult(raw)
	if !ok {
		repaired, repairErr := s.llm.Complete(ctx, model,
			"Fix the JSON. Return exactly the requested object with no markdown or additional prose.",
			"Invalid output:\n"+limit(raw, 5000)+"\n\nRequired shape:\n"+resultSchema,
		)
		if repairErr != nil {
			if IsPermanent(repairErr) {
				return s.FailAndRefund(ctx, replyGenerationID, repairErr)
			}
			return repairErr
		}
		result, ok = decodeResult(repaired)
		if !ok {
			return s.FailAndRefund(ctx, replyGenerationID, ErrInvalidOutput)
		}
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return s.FailAndRefund(ctx, replyGenerationID, ErrInvalidOutput)
	}
	_, err = s.repository.CompleteDocumentReplyGeneration(ctx, db.CompleteDocumentReplyGenerationParams{
		ID:     replyGenerationID,
		Result: resultJSON,
	})
	return err
}

func (s *Service) FailAndRefund(ctx context.Context, replyGenerationID uuid.UUID, cause error) error {
	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin terminal generation failure transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	generation, err := s.repository.FailActiveDocumentReplyGenerationTx(ctx, tx, db.FailActiveDocumentReplyGenerationParams{
		ID:           replyGenerationID,
		ErrorMessage: pgtype.Text{String: sanitizedGenerationError(cause), Valid: true},
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("mark generation terminally failed: %w", err)
	}

	if _, err := s.points.RefundTx(ctx, tx, users_service.RefundTxInput{
		UserID:            generation.UserID,
		Amount:            GenerationCost,
		Description:       "Papierbrief generation refund",
		RelatedEntityID:   generation.ID,
		RelatedEntityType: db.RelatedEntityTypePaperLetterGeneration,
	}); err != nil && !errors.Is(err, users_service.ErrBuddyPointSpendNotFound) {
		return fmt.Errorf("refund terminally failed generation: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit terminal generation failure: %w", err)
	}
	return cause
}

func sanitizedGenerationError(cause error) string {
	switch {
	case errors.Is(cause, ErrInvalidInput):
		return "invalid_generation_input"
	case errors.Is(cause, ErrExtractionNotReady), errors.Is(cause, ErrAnalysisNotReady):
		return "source_not_ready"
	case errors.Is(cause, llm_service.ErrUnsupportedModel):
		return "unsupported_model"
	case errors.Is(cause, ErrInvalidOutput):
		return "invalid_generation_output"
	default:
		return "generation_failed"
	}
}

func IsPermanent(err error) bool {
	switch {
	case errors.Is(err, ErrGenerationNotFound),
		errors.Is(err, ErrInvalidInput),
		errors.Is(err, ErrExtractionNotReady),
		errors.Is(err, ErrAnalysisNotReady),
		errors.Is(err, ErrInvalidOutput),
		errors.Is(err, llm_service.ErrMissingAPIKey),
		errors.Is(err, llm_service.ErrUnsupportedModel):
		return true
	}

	var retryable interface{ Retryable() bool }
	return errors.As(err, &retryable) && !retryable.Retryable()
}

const resultSchema = `{
  "germanLetter": "letter body in German without sender/recipient addresses, date, or subject",
  "userLanguageLetter": "the same letter translated into the requested language",
  "explanation": "short explanation in the requested language",
  "subject": "short German subject line",
  "tone": "neutral | clarification | objection | information"
}`

func buildPrompt(mode db.DocumentReplyGenerationMode, language string, input GenerationInput, sourceText string, analysisMetadata []byte) (string, string) {
	if mode == db.DocumentReplyGenerationModeImprove {
		return buildImprovePrompt(language, input, sourceText, analysisMetadata)
	}

	system := "You write formal responses to official letters and bureaucratic institutions in Germany. " +
		"The outgoing letter must be in German, polite, factual, and must not invent facts. " +
		"Return only the letter body; omit sender and recipient addresses, the date, and the subject because the PDF layout adds them. " +
		"Treat a submitted outgoing sender and recipient as authoritative and never swap their roles. " +
		"When no submitted parties are provided, the receiver of the original letter is the sender of the reply, " +
		"and the sender of the original letter is the recipient of the reply. " +
		"Use source OCR and analysis only as supporting context. Never replace submitted facts with conflicting source data. " +
		"If a needed fact is missing, state the gap or ask for clarification instead of inventing it. " +
		"Return strict JSON only."

	instruction := strings.TrimSpace(input.Instruction)
	if instruction == "" {
		if strings.TrimSpace(sourceText) == "" && len(analysisMetadata) == 0 {
			instruction = "Draft a formal letter using only the submitted context."
		} else {
			instruction = "Draft a formal response based on the original letter and its analysis."
		}
	}
	user := fmt.Sprintf(
		"Translation and explanation language: %s\n\nSubmitted outgoing letter context:\n%s\n\nUser instruction:\n%s\n\nAnalysis of the original letter:\n%s\n\nOriginal OCR text:\n%s\n\nRequired output:\n%s",
		languageLabel(language),
		generationContext(input),
		limit(instruction, 2500),
		optionalContext(limit(string(analysisMetadata), 24000)),
		optionalContext(limit(sourceText, 24000)),
		resultSchema,
	)
	return system, user
}

func buildImprovePrompt(language string, input GenerationInput, sourceText string, analysisMetadata []byte) (string, string) {
	system := "You improve an existing German letter addressed to an official or bureaucratic institution in Germany. " +
		"Preserve every fact, request, commitment, and the writer's intent unless the submitted feedback explicitly asks for a change. " +
		"Return only the letter body; omit sender and recipient addresses, the date, and the subject because the PDF layout adds them. " +
		"Treat the submitted outgoing sender and recipient as authoritative and never swap their roles. " +
		"Do not invent names, dates, reference numbers, events, legal claims, or other facts. " +
		"Use source OCR and analysis only as supporting context and never let them silently change the submitted letter. " +
		"If the feedback requires missing information, state the gap instead of inventing it. " +
		"Return strict JSON only."

	user := fmt.Sprintf(
		"Translation and explanation language: %s\n\nSubmitted outgoing letter context:\n%s\n\nCurrent German letter to improve:\n%s\n\nRequested improvements:\n%s\n\nAnalysis of the original incoming letter (supporting context only):\n%s\n\nOriginal incoming OCR text (supporting context only):\n%s\n\nRequired output:\n%s",
		languageLabel(language),
		generationContext(input),
		limit(input.CurrentGermanLetter, 24000),
		limit(input.Feedback, 2500),
		optionalContext(limit(string(analysisMetadata), 24000)),
		optionalContext(limit(sourceText, 24000)),
		resultSchema,
	)
	return system, user
}

func hasSubmittedContext(input GenerationInput) bool {
	return hasPartyContext(input.Sender) || hasPartyContext(input.Recipient) ||
		strings.TrimSpace(input.Category) != "" || strings.TrimSpace(input.Subcategory) != "" ||
		strings.TrimSpace(input.Subject) != "" || strings.TrimSpace(input.Reference) != "" ||
		strings.TrimSpace(input.Instruction) != ""
}

func hasPartyContext(party *PartyInput) bool {
	return party != nil && (strings.TrimSpace(party.Name) != "" || strings.TrimSpace(party.Address) != "")
}

func generationContext(input GenerationInput) string {
	return fmt.Sprintf(
		"Outgoing sender:\n%s\n\nOutgoing recipient:\n%s\n\nCategory: %s\nSubcategory: %s\nSubject: %s\nReference: %s\nUI locale: %s",
		partyContext(input.Sender),
		partyContext(input.Recipient),
		optionalContext(limit(input.Category, 100)),
		optionalContext(limit(input.Subcategory, 100)),
		optionalContext(limit(input.Subject, 500)),
		optionalContext(limit(input.Reference, 200)),
		optionalContext(limit(input.Locale, 10)),
	)
}

func partyContext(party *PartyInput) string {
	if party == nil {
		return "(not provided)"
	}
	return fmt.Sprintf(
		"Name: %s\nAddress: %s",
		optionalContext(limit(party.Name, 200)),
		optionalContext(limit(party.Address, 1000)),
	)
}

func optionalContext(value string) string {
	if strings.TrimSpace(value) == "" {
		return "(not provided)"
	}
	return value
}

func decodeResult(raw string) (Result, bool) {
	var result Result
	value := strings.TrimSpace(raw)
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		match := regexp.MustCompile(`\{[\s\S]*\}`).FindString(value)
		if match == "" || json.Unmarshal([]byte(match), &result) != nil {
			return Result{}, false
		}
	}

	result.GermanLetter = strings.TrimSpace(result.GermanLetter)
	result.UserLanguageLetter = strings.TrimSpace(result.UserLanguageLetter)
	result.Explanation = strings.TrimSpace(result.Explanation)
	result.Subject = strings.TrimSpace(result.Subject)
	result.Tone = strings.ToLower(strings.TrimSpace(result.Tone))
	if result.GermanLetter == "" {
		return Result{}, false
	}
	if result.UserLanguageLetter == "" {
		result.UserLanguageLetter = result.GermanLetter
	}
	switch result.Tone {
	case "neutral", "clarification", "objection", "information":
	default:
		result.Tone = "neutral"
	}
	return result, true
}

func languageLabel(code string) string {
	switch code {
	case "en":
		return "English"
	case "ru":
		return "Russian"
	case "uk", "ua":
		return "Ukrainian"
	case "pl":
		return "Polish"
	case "ar":
		return "Arabic"
	default:
		return "German"
	}
}

func limit(value string, maxRunes int) string {
	runes := []rune(strings.TrimSpace(value))
	if len(runes) <= maxRunes {
		return string(runes)
	}
	return string(runes[:maxRunes])
}
