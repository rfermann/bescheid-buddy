package analysis_service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"regexp"
	"strings"

	"github.com/abadojack/whatlanggo"
	"github.com/bored-bots/bored-template/internal/repository/cases_repository"
	"github.com/bored-bots/bored-template/internal/services/llm_service"
	"github.com/bored-bots/bored-template/internal/store/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const analysisMaxChars = 32000

var (
	ErrExtractionNotReady    = errors.New("document extraction is not completed")
	ErrNoTextToAnalyze       = errors.New("no extracted text to analyze")
	ErrInvalidAnalysisOutput = errors.New("analysis model did not return valid output")
)

type Service struct {
	repository cases_repository.Repository
	llm        *llm_service.Service
}

func New(repository cases_repository.Repository, llm *llm_service.Service) *Service {
	return &Service{repository: repository, llm: llm}
}

func (s *Service) AnalyzeDocument(ctx context.Context, documentID uuid.UUID, modelStr, language string) error {
	model, err := llm_service.ParseModel(modelStr)
	if err != nil {
		return s.fail(ctx, documentID, err)
	}

	if _, err := s.repository.MarkDocumentAnalysisProcessing(ctx, db.MarkDocumentAnalysisProcessingParams{
		DocumentID: documentID,
	}); err != nil {
		return err
	}

	text, err := s.extractedText(ctx, documentID)
	if err != nil {
		return s.fail(ctx, documentID, err)
	}

	analysisText, metadata, detectedLang, title, err := s.analyze(ctx, model, text, pickLang(language))
	if err != nil {
		if IsPermanent(err) {
			return s.fail(ctx, documentID, err)
		}
		return err
	}

	var caseID uuid.UUID
	if title != "" {
		document, err := s.repository.GetCaseDocumentByID(ctx, db.GetCaseDocumentByIDParams{ID: documentID})
		if err != nil {
			return err
		}
		initialDocument, err := s.repository.GetInitialCaseDocument(ctx, db.GetInitialCaseDocumentParams{CaseID: document.CaseID})
		if err != nil {
			return err
		}
		if initialDocument.ID == documentID {
			caseID = document.CaseID
		}
	}

	tx, err := s.repository.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	txRepo := s.repository.WithTx(tx)
	if _, err = txRepo.CompleteDocumentAnalysis(ctx, db.CompleteDocumentAnalysisParams{
		DocumentID:       documentID,
		AnalysisText:     pgtype.Text{String: analysisText, Valid: true},
		AnalysisMetadata: metadata,
		Provider:         pgtype.Text{String: model.Provider(), Valid: true},
		DetectedLanguage: pgtype.Text{String: detectedLang, Valid: detectedLang != ""},
	}); err != nil {
		return err
	}
	if caseID != uuid.Nil {
		if err := txRepo.UpdateCaseTitle(ctx, db.UpdateCaseTitleParams{
			ID:    caseID,
			Title: pgtype.Text{String: title, Valid: true},
		}); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (s *Service) extractedText(ctx context.Context, documentID uuid.UUID) (string, error) {
	row, err := s.repository.GetDocumentExtractionTextByDocumentID(ctx, db.GetDocumentExtractionTextByDocumentIDParams{
		DocumentID: documentID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrExtractionNotReady
		}
		return "", err
	}
	if row.Status != db.DocumentAnalysisStatusCompleted {
		return "", ErrExtractionNotReady
	}
	text := strings.TrimSpace(row.Text.String)
	if text == "" {
		return "", ErrNoTextToAnalyze
	}
	return text, nil
}

func (s *Service) analyze(ctx context.Context, model llm_service.Model, text, langCode string) (string, []byte, string, string, error) {
	label := langLabel(langCode)
	prepared := prepareText(text, analysisMaxChars)
	system, user := buildPrompt(prepared, label)

	raw, err := s.llm.Complete(ctx, model, system, user)
	if err != nil {
		return "", nil, "", "", err
	}

	decoded, ok := decodeJSON(raw)
	if !ok || !hasAnalysisContent(decoded) {
		repaired, rerr := s.llm.Complete(ctx, model, repairSystemPrompt, repairUserPrompt(raw))
		if rerr != nil {
			return "", nil, "", "", rerr
		}
		decoded, ok = decodeJSON(repaired)
		if !ok || !hasAnalysisContent(decoded) {
			return "", nil, "", "", ErrInvalidAnalysisOutput
		}
	}

	detectedLang := s.enforceTranslationLanguage(ctx, model, langCode, label, decoded)

	metadata, err := json.Marshal(decoded)
	if err != nil {
		return "", nil, "", "", ErrInvalidAnalysisOutput
	}
	return analysisText(decoded), metadata, detectedLang, caseTitle(decoded), nil
}

var analysisLangWhitelist = map[whatlanggo.Lang]bool{
	whatlanggo.Deu: true,
	whatlanggo.Eng: true,
	whatlanggo.Rus: true,
	whatlanggo.Ukr: true,
	whatlanggo.Pol: true,
	whatlanggo.Arb: true,
}

const (
	minTranslationBytes = 40
	minLangConfidence   = 0.65
)

func (s *Service) enforceTranslationLanguage(ctx context.Context, model llm_service.Model, langCode, label string, decoded map[string]any) string {
	part1 := strings.TrimSpace(stringField(decoded, "part_1_full_translation"))
	if len(part1) < minTranslationBytes {
		return ""
	}

	info := whatlanggo.DetectWithOptions(part1, whatlanggo.Options{Whitelist: analysisLangWhitelist})
	detected := info.Lang.Iso6391()
	if detected == "" || info.Confidence < minLangConfidence {
		return ""
	}
	if detected == iso6391(langCode) {
		return detected
	}

	retranslated, err := s.retranslate(ctx, model, label, part1)
	if err != nil || retranslated == "" {
		slog.WarnContext(ctx, "analysis language guard: re-translation failed",
			"model", string(model), "target", langCode, "detected", detected, "error", err)
		return detected
	}

	decoded["part_1_full_translation"] = retranslated
	slog.InfoContext(ctx, "analysis language guard: re-translated PART 1",
		"model", string(model), "target", langCode, "detected", detected, "confidence", info.Confidence)
	return detected
}

func (s *Service) retranslate(ctx context.Context, model llm_service.Model, label, text string) (string, error) {
	system := "You are a translator. Translate the user's text into " + label + ". " +
		"Respond with strict JSON only: {\"translation\": \"<the full text translated into " + label + ">\"}. " +
		"No markdown, no extra keys, no commentary."

	raw, err := s.llm.Complete(ctx, model, system, text)
	if err != nil {
		return "", err
	}
	m, ok := decodeJSON(raw)
	if !ok {
		return "", nil
	}
	return strings.TrimSpace(stringField(m, "translation")), nil
}

func iso6391(code string) string {
	switch pickLang(code) {
	case "ua", "uk":
		return "uk"
	default:
		return pickLang(code)
	}
}

func (s *Service) fail(ctx context.Context, documentID uuid.UUID, cause error) error {
	if _, err := s.repository.FailDocumentAnalysis(ctx, db.FailDocumentAnalysisParams{
		DocumentID:   documentID,
		ErrorMessage: pgtype.Text{String: cause.Error(), Valid: true},
	}); err != nil {
		return err
	}
	return cause
}

func IsPermanent(err error) bool {
	switch {
	case errors.Is(err, ErrExtractionNotReady),
		errors.Is(err, ErrNoTextToAnalyze),
		errors.Is(err, ErrInvalidAnalysisOutput),
		errors.Is(err, llm_service.ErrMissingAPIKey),
		errors.Is(err, llm_service.ErrUnsupportedModel):
		return true
	}

	var retryabler interface{ Retryable() bool }
	if errors.As(err, &retryabler) {
		return !retryabler.Retryable()
	}

	return false
}

const schemaExample = `{
  "receiver": {
    "name": "",
    "street": "",
    "house_number": "",
    "postal_code": "",
    "city": "",
    "state": "",
    "country": ""
  },
  "sender": {
    "name": "",
    "street": "",
    "house_number": "",
    "postal_code": "",
    "city": "",
    "state": "",
    "country": "",
    "category": "unknown"
  },
  "letter_type": "unknown",
  "topic": "",
  "subject": {
    "localized": "",
    "original": ""
  },
  "dates": { "letter_date": "", "deadline_date": "" },
  "importance": "medium",
  "part_1_full_translation": "",
  "part_2_ai_analysis": {
    "summary": "",
    "documents_requested": [],
    "deadlines": [],
    "recommended_next_steps": [],
    "possible_consequences": []
  }
}`

func pickLang(code string) string {
	c := strings.ToLower(strings.TrimSpace(code))
	switch c {
	case "de", "en", "ru", "uk", "ua", "pl", "ar":
		return c
	default:
		return "de"
	}
}

func langLabel(code string) string {
	switch pickLang(code) {
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

func buildPrompt(text, languageLabel string) (system, user string) {
	system = "You are an expert AI assistant specialized in analyzing official letters, especially bureaucratic letters from government authorities, insurance companies, banks and courts.\n\n" +
		"Your goal is to help the user clearly understand the letter they received.\n\n" +
		"OUTPUT LANGUAGE — THIS IS CRITICAL\n" +
		"Every field you output, in BOTH parts, MUST be written in " + languageLabel + " and only in " + languageLabel + ", except subject.original.\n\n" +
		"The output MUST always contain TWO CLEAR PARTS:\n" +
		"PART 1 (part_1_full_translation): a COMPLETE TRANSLATION of the letter INTO " + languageLabel + ". " +
		"The letter itself is usually written in German (or another language). You MUST translate all of its content into " + languageLabel + ". " +
		"Do NOT copy or transcribe the original wording — never leave the German (or source-language) text in this field. " +
		"Translate it into " + languageLabel + ". The only case where the text may stay unchanged is when the letter is already written in " + languageLabel + ".\n" +
		"PART 2 (part_2_ai_analysis): your analysis and explanation of the letter, also written entirely in " + languageLabel + ".\n\n" +
		"IMPORTANT RULES\n" +
		"- Do NOT invent facts not present in the letter.\n" +
		"- If something is unclear or missing, explicitly say it is unclear.\n" +
		"- subject.original: copy the letter's explicitly printed subject exactly in its source language. Do not infer a subject from the body.\n" +
		"- subject.localized: translate subject.original faithfully into " + languageLabel + ". If the letter has no explicit subject, return an empty string for both subject fields.\n" +
		"- Extract sender, topic, reference numbers and deadlines if possible.\n" +
		"- Use clear, simple language.\n" +
		"- Do not give legal advice - only explanation and possible next steps.\n\n" +
		"Return strict JSON only. No markdown, no code fences, no prose outside JSON. Use exactly the keys from the output format."

	user = "Target output language: " + languageLabel + ". Every field below except subject.original must be written in " + languageLabel + ".\n\n" +
		"OCR text of the letter:\n" + text + "\n\n" +
		"OUTPUT FORMAT (STRICT JSON) — remember: part_1_full_translation must be the letter TRANSLATED INTO " + languageLabel + ", not the original text:\n" + schemaExample
	return system, user
}

const repairSystemPrompt = "Fix your JSON. Return valid JSON only, matching the required schema exactly. No markdown, no code fences, no prose outside JSON."

func repairUserPrompt(previous string) string {
	if len(previous) > 4000 {
		previous = previous[:4000]
	}
	return "Previous invalid output:\n" + previous + "\n\nReturn corrected JSON only, using this schema:\n" + schemaExample
}

var reFirstObject = regexp.MustCompile(`(?s)\{.*\}`)

func decodeJSON(raw string) (map[string]any, bool) {
	trimmed := strings.TrimSpace(raw)
	var m map[string]any
	if err := json.Unmarshal([]byte(trimmed), &m); err == nil {
		return m, true
	}
	if match := reFirstObject.FindString(trimmed); match != "" {
		if err := json.Unmarshal([]byte(match), &m); err == nil {
			return m, true
		}
	}
	return nil, false
}

func hasAnalysisContent(m map[string]any) bool {
	return analysisText(m) != ""
}

func analysisText(m map[string]any) string {
	if t := strings.TrimSpace(stringField(m, "part_1_full_translation")); t != "" {
		return t
	}
	if part2, ok := m["part_2_ai_analysis"].(map[string]any); ok {
		return strings.TrimSpace(stringField(part2, "summary"))
	}
	return ""
}

func stringField(m map[string]any, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func caseTitle(m map[string]any) string {
	subject, ok := m["subject"].(map[string]any)
	if !ok {
		return ""
	}

	localized := strings.Join(strings.Fields(stringField(subject, "localized")), " ")
	original := strings.Join(strings.Fields(stringField(subject, "original")), " ")
	if localized == "" {
		return truncateTitle(original, 255)
	}
	if original == "" || strings.EqualFold(localized, original) {
		return truncateTitle(localized, 255)
	}

	localized = truncateTitle(localized, 255)
	remaining := 255 - len([]rune(localized)) - len([]rune(" / "))
	if remaining <= 0 {
		return localized
	}
	return localized + " / " + truncateTitle(original, remaining)
}

func truncateTitle(title string, maxChars int) string {
	runes := []rune(title)
	if len(runes) > maxChars {
		return string(runes[:maxChars])
	}
	return title
}

func prepareText(text string, maxChars int) string {
	runes := []rune(text)
	if len(runes) <= maxChars {
		return text
	}
	headLen := maxChars * 7 / 10
	tailLen := maxChars - headLen
	head := string(runes[:headLen])
	tail := string(runes[len(runes)-tailLen:])
	return head + "\n\n[... OCR TEXT TRUNCATED ...]\n\n" + tail
}
