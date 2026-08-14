package reply_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/bored-bots/bored-template/internal/services/llm_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrGenerationNotFound = errors.New("reply generation was not found")
	ErrExtractionNotReady = errors.New("document extraction is not completed")
	ErrAnalysisNotReady   = errors.New("document analysis is not completed")
	ErrInvalidOutput      = errors.New("reply model did not return valid output")
)

type Repository interface {
	GetDocumentReplyGenerationByID(ctx context.Context, params db.GetDocumentReplyGenerationByIDParams) (*db.DocumentReplyGeneration, error)
	MarkDocumentReplyGenerationProcessing(ctx context.Context, params db.MarkDocumentReplyGenerationProcessingParams) (*db.DocumentReplyGeneration, error)
	CompleteDocumentReplyGeneration(ctx context.Context, params db.CompleteDocumentReplyGenerationParams) (*db.DocumentReplyGeneration, error)
	FailDocumentReplyGeneration(ctx context.Context, params db.FailDocumentReplyGenerationParams) (*db.DocumentReplyGeneration, error)
	GetDocumentExtractionTextByDocumentID(ctx context.Context, params db.GetDocumentExtractionTextByDocumentIDParams) (*db.GetDocumentExtractionTextByDocumentIDRow, error)
	GetDocumentAnalysisByDocumentID(ctx context.Context, params db.GetDocumentAnalysisByDocumentIDParams) (*db.DocumentAnalysisResult, error)
}

type LLM interface {
	Complete(ctx context.Context, model llm_service.Model, system, user string) (string, error)
}

type Service struct {
	repository Repository
	llm        LLM
}

type Result struct {
	GermanLetter       string `json:"germanLetter"`
	UserLanguageLetter string `json:"userLanguageLetter"`
	Explanation        string `json:"explanation"`
	Subject            string `json:"subject"`
	Tone               string `json:"tone"`
}

func New(repository Repository, llm LLM) *Service {
	return &Service{repository: repository, llm: llm}
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

	model, err := llm_service.ParseModel(generation.Model)
	if err != nil {
		return s.fail(ctx, replyGenerationID, err)
	}
	if _, err := s.repository.MarkDocumentReplyGenerationProcessing(ctx, db.MarkDocumentReplyGenerationProcessingParams{ID: replyGenerationID}); err != nil {
		return err
	}

	extraction, err := s.repository.GetDocumentExtractionTextByDocumentID(ctx, db.GetDocumentExtractionTextByDocumentIDParams{DocumentID: generation.DocumentID})
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && extraction.Status != db.DocumentAnalysisStatusCompleted) {
		return s.fail(ctx, replyGenerationID, ErrExtractionNotReady)
	}
	if err != nil {
		return err
	}

	analysis, err := s.repository.GetDocumentAnalysisByDocumentID(ctx, db.GetDocumentAnalysisByDocumentIDParams{DocumentID: generation.DocumentID})
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && analysis.Status != db.DocumentAnalysisStatusCompleted) {
		return s.fail(ctx, replyGenerationID, ErrAnalysisNotReady)
	}
	if err != nil {
		return err
	}

	system, user := buildPrompt(
		generation.Language,
		generation.Instruction,
		extraction.Text.String,
		analysis.AnalysisMetadata,
	)
	raw, err := s.llm.Complete(ctx, model, system, user)
	if err != nil {
		if IsPermanent(err) {
			return s.fail(ctx, replyGenerationID, err)
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
				return s.fail(ctx, replyGenerationID, repairErr)
			}
			return repairErr
		}
		result, ok = decodeResult(repaired)
		if !ok {
			return s.fail(ctx, replyGenerationID, ErrInvalidOutput)
		}
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return s.fail(ctx, replyGenerationID, ErrInvalidOutput)
	}
	_, err = s.repository.CompleteDocumentReplyGeneration(ctx, db.CompleteDocumentReplyGenerationParams{
		ID:     replyGenerationID,
		Result: resultJSON,
	})
	return err
}

func (s *Service) fail(ctx context.Context, replyGenerationID uuid.UUID, cause error) error {
	if _, err := s.repository.FailDocumentReplyGeneration(ctx, db.FailDocumentReplyGenerationParams{
		ID:           replyGenerationID,
		ErrorMessage: pgtype.Text{String: cause.Error(), Valid: true},
	}); err != nil {
		return err
	}
	return cause
}

func IsPermanent(err error) bool {
	switch {
	case errors.Is(err, ErrGenerationNotFound),
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
  "germanLetter": "full formal response letter in German",
  "userLanguageLetter": "the same letter translated into the requested language",
  "explanation": "short explanation in the requested language",
  "subject": "short German subject line",
  "tone": "neutral | clarification | objection | information"
}`

func buildPrompt(language, instruction, sourceText string, analysisMetadata []byte) (string, string) {
	system := "You write formal responses to official letters and bureaucratic institutions in Germany. " +
		"The outgoing letter must be in German, polite, factual, and must not invent facts. " +
		"The receiver of the original letter is the sender of the reply. " +
		"The sender of the original letter is the recipient of the reply. " +
		"If information is missing, ask for clarification. Return strict JSON only."

	if strings.TrimSpace(instruction) == "" {
		instruction = "Draft a formal response based on the original letter and its analysis."
	}
	user := fmt.Sprintf(
		"Translation and explanation language: %s\n\nUser instruction:\n%s\n\nAnalysis of the original letter:\n%s\n\nOriginal OCR text:\n%s\n\nRequired output:\n%s",
		languageLabel(language),
		limit(instruction, 2500),
		limit(string(analysisMetadata), 24000),
		limit(sourceText, 24000),
		resultSchema,
	)
	return system, user
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
