package llm_service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Model string

const (
	ModelGPT56Luna     Model = "gpt-5.6-luna"
	ModelGeminiFlash37 Model = "gemini-3.7-flash"
	ModelHaiku45       Model = "claude-haiku-4-5"
)

type ModelInfo struct {
	ID       Model
	Name     string
	Provider string
}

var models = []ModelInfo{
	{ID: ModelGPT56Luna, Name: "GPT-5.6 Luna", Provider: "openai"},
	{ID: ModelGeminiFlash37, Name: "Gemini 3.7 Flash", Provider: "google"},
	{ID: ModelHaiku45, Name: "Claude Haiku 4.5", Provider: "anthropic"},
}

const maxOutputTokens = 2048

var (
	ErrUnsupportedModel = errors.New("unsupported model")
	ErrMissingAPIKey    = errors.New("no api key configured for the selected model's provider")
)

func ParseModel(s string) (Model, error) {
	for _, model := range models {
		if model.ID == Model(s) {
			return model.ID, nil
		}
	}
	return "", fmt.Errorf("%w: %s", ErrUnsupportedModel, s)
}

func Models() []ModelInfo {
	return append([]ModelInfo(nil), models...)
}

func (m Model) Provider() string {
	for _, model := range models {
		if model.ID == m {
			return model.Provider
		}
	}
	return "unknown"
}

type Config struct {
	OpenAIAPIKey    string
	GeminiAPIKey    string
	AnthropicAPIKey string
}

type Service struct {
	cfg        Config
	httpClient *http.Client
}

func New(cfg Config) *Service {
	return &Service{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 120 * time.Second},
	}
}

type APIError struct {
	Provider   string
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s api status %d: %s", e.Provider, e.StatusCode, truncate(e.Body, 400))
}

func (e *APIError) Retryable() bool {
	return e.StatusCode == http.StatusTooManyRequests || e.StatusCode >= 500
}

func (s *Service) Complete(ctx context.Context, model Model, system, user string) (string, error) {
	switch model.Provider() {
	case "openai":
		return s.completeOpenAI(ctx, model, system, user)
	case "google":
		return s.completeGemini(ctx, model, system, user)
	case "anthropic":
		return s.completeAnthropic(ctx, model, system, user)
	default:
		return "", fmt.Errorf("%w: %s", ErrUnsupportedModel, model)
	}
}

func (s *Service) completeOpenAI(ctx context.Context, model Model, system, user string) (string, error) {
	if s.cfg.OpenAIAPIKey == "" {
		return "", ErrMissingAPIKey
	}

	payload := openAIRequest{
		Model:           string(model),
		ReasoningEffort: "none",
		Messages: []openAIMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
		ResponseFormat: &openAIResponseFormat{Type: "json_object"},
	}

	req, err := s.newJSONRequest(ctx, "https://api.openai.com/v1/chat/completions", payload)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.cfg.OpenAIAPIKey)

	var resp openAIResponse
	if err := s.do(ctx, "openai", req, &resp); err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}
	return resp.Choices[0].Message.Content, nil
}

func (s *Service) completeGemini(ctx context.Context, model Model, system, user string) (string, error) {
	if s.cfg.GeminiAPIKey == "" {
		return "", ErrMissingAPIKey
	}

	payload := geminiRequest{
		SystemInstruction: &geminiContent{Parts: []geminiPart{{Text: system}}},
		Contents:          []geminiContent{{Role: "user", Parts: []geminiPart{{Text: user}}}},
		GenerationConfig:  &geminiGenerationConfig{ResponseMimeType: "application/json"},
	}

	endpoint := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		string(model), url.QueryEscape(s.cfg.GeminiAPIKey),
	)
	req, err := s.newJSONRequest(ctx, endpoint, payload)
	if err != nil {
		return "", err
	}

	var resp geminiResponse
	if err := s.do(ctx, "google", req, &resp); err != nil {
		return "", err
	}
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned no content")
	}
	return resp.Candidates[0].Content.Parts[0].Text, nil
}

func (s *Service) completeAnthropic(ctx context.Context, model Model, system, user string) (string, error) {
	if s.cfg.AnthropicAPIKey == "" {
		return "", ErrMissingAPIKey
	}

	payload := anthropicRequest{
		Model:     string(model),
		MaxTokens: maxOutputTokens,
		System:    system,
		Messages:  []anthropicMessage{{Role: "user", Content: user}},
	}

	req, err := s.newJSONRequest(ctx, "https://api.anthropic.com/v1/messages", payload)
	if err != nil {
		return "", err
	}
	req.Header.Set("x-api-key", s.cfg.AnthropicAPIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	var resp anthropicResponse
	if err := s.do(ctx, "anthropic", req, &resp); err != nil {
		return "", err
	}
	for _, block := range resp.Content {
		if block.Type == "text" && block.Text != "" {
			return block.Text, nil
		}
	}
	return "", fmt.Errorf("anthropic returned no text content")
}

func (s *Service) newJSONRequest(ctx context.Context, endpoint string, payload any) (*http.Request, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (s *Service) do(ctx context.Context, provider string, req *http.Request, out any) error {
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{Provider: provider, StatusCode: resp.StatusCode, Body: string(respBody)}
	}
	return json.Unmarshal(respBody, out)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

type openAIResponseFormat struct {
	Type string `json:"type"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIRequest struct {
	Model           string                `json:"model"`
	ReasoningEffort string                `json:"reasoning_effort,omitempty"`
	Messages        []openAIMessage       `json:"messages"`
	ResponseFormat  *openAIResponseFormat `json:"response_format,omitempty"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiGenerationConfig struct {
	ResponseMimeType string `json:"responseMimeType,omitempty"`
}

type geminiRequest struct {
	SystemInstruction *geminiContent          `json:"systemInstruction,omitempty"`
	Contents          []geminiContent         `json:"contents"`
	GenerationConfig  *geminiGenerationConfig `json:"generationConfig,omitempty"`
}

type geminiResponse struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}
