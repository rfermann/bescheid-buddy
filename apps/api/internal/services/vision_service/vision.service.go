package vision_service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultBaseURL = "https://vision.googleapis.com/v1"

type Service struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func New(apiKey string) *Service {
	return &Service{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("vision api status %d: %s", e.StatusCode, truncate(e.Body, 400))
}

func (e *APIError) Retryable() bool {
	return e.StatusCode == http.StatusTooManyRequests || e.StatusCode >= 500
}

func (s *Service) DetectImageText(ctx context.Context, data []byte) (string, error) {
	payload := imagesAnnotatePayload{
		Requests: []imageAnnotateRequest{{
			Image:    imageContent{Content: base64.StdEncoding.EncodeToString(data)},
			Features: []annotateFeature{{Type: "DOCUMENT_TEXT_DETECTION", MaxResults: 1}},
		}},
	}

	var resp imagesAnnotateResponse
	if err := s.do(ctx, "images:annotate", payload, &resp); err != nil {
		return "", err
	}
	if len(resp.Responses) == 0 {
		return "", nil
	}
	if resp.Responses[0].Error != nil {
		return "", fmt.Errorf("vision image annotation failed: %s", resp.Responses[0].Error.Message)
	}
	return extractText(resp.Responses[0]), nil
}

func (s *Service) DetectPDFText(ctx context.Context, data []byte, pages []int) (string, error) {
	payload := filesAnnotatePayload{
		Requests: []fileAnnotateRequest{{
			InputConfig: fileInputConfig{
				Content:  base64.StdEncoding.EncodeToString(data),
				MimeType: "application/pdf",
			},
			Features: []annotateFeature{{Type: "DOCUMENT_TEXT_DETECTION"}},
			Pages:    pages,
		}},
	}

	var resp filesAnnotateResponse
	if err := s.do(ctx, "files:annotate", payload, &resp); err != nil {
		return "", err
	}
	if len(resp.Responses) == 0 {
		return "", nil
	}
	if resp.Responses[0].Error != nil {
		return "", fmt.Errorf("vision file annotation failed: %s", resp.Responses[0].Error.Message)
	}

	var b strings.Builder
	for _, page := range resp.Responses[0].Responses {
		t := extractText(page)
		if t == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString(t)
	}
	return b.String(), nil
}

func (s *Service) do(ctx context.Context, endpoint string, payload any, out any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	endpointURL := fmt.Sprintf("%s/%s?key=%s", s.baseURL, endpoint, url.QueryEscape(s.apiKey))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpointURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

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
		return &APIError{StatusCode: resp.StatusCode, Body: string(respBody)}
	}

	return json.Unmarshal(respBody, out)
}

func extractText(r annotateImageResponse) string {
	if t := strings.TrimSpace(r.FullTextAnnotation.Text); t != "" {
		return t
	}
	if len(r.TextAnnotations) > 0 {
		return strings.TrimSpace(r.TextAnnotations[0].Description)
	}
	return ""
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

type annotateFeature struct {
	Type       string `json:"type"`
	MaxResults int    `json:"maxResults,omitempty"`
}

type imageContent struct {
	Content string `json:"content"`
}

type imageAnnotateRequest struct {
	Image    imageContent      `json:"image"`
	Features []annotateFeature `json:"features"`
}

type imagesAnnotatePayload struct {
	Requests []imageAnnotateRequest `json:"requests"`
}

type fileInputConfig struct {
	Content  string `json:"content"`
	MimeType string `json:"mimeType"`
}

type fileAnnotateRequest struct {
	InputConfig fileInputConfig   `json:"inputConfig"`
	Features    []annotateFeature `json:"features"`
	Pages       []int             `json:"pages,omitempty"`
}

type filesAnnotatePayload struct {
	Requests []fileAnnotateRequest `json:"requests"`
}

type visionStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type fullTextAnnotation struct {
	Text string `json:"text"`
}

type textAnnotation struct {
	Description string `json:"description"`
}

type annotateImageResponse struct {
	FullTextAnnotation fullTextAnnotation `json:"fullTextAnnotation"`
	TextAnnotations    []textAnnotation   `json:"textAnnotations"`
	Error              *visionStatus      `json:"error,omitempty"`
}

type imagesAnnotateResponse struct {
	Responses []annotateImageResponse `json:"responses"`
}

type fileAnnotateResponse struct {
	Responses []annotateImageResponse `json:"responses"`
	Error     *visionStatus           `json:"error,omitempty"`
}

type filesAnnotateResponse struct {
	Responses []fileAnnotateResponse `json:"responses"`
}
