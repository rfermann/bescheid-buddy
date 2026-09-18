package pingen_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bored-bots/bored-template/internal/config"
)

type Service struct {
	config config.PingenConfig
	client *http.Client
}

type Error struct {
	operation string
	status    int
	code      string
	temporary bool
}

func (e *Error) Error() string {
	if e.status != 0 {
		return fmt.Sprintf("Pingen %s failed with status %d", e.operation, e.status)
	}
	return "Pingen " + e.operation + " failed"
}

func (e *Error) Temporary() bool { return e.temporary }

func (e *Error) SafeMessage() string {
	if e.code != "" {
		return fmt.Sprintf("pingen_status_%d_%s", e.status, e.code)
	}
	return fmt.Sprintf("pingen_status_%d", e.status)
}

func New(cfg config.PingenConfig) *Service {
	return &Service{config: cfg, client: &http.Client{Timeout: 20 * time.Second}}
}

func (s *Service) Submit(ctx context.Context, pdf []byte, filename, idempotencyKey, product, addressPosition, printMode, printSpectrum string) (string, error) {
	token, err := s.token(ctx)
	if err != nil {
		return "", err
	}

	var upload struct {
		Data struct {
			Attributes struct {
				URL       string `json:"url"`
				Signature string `json:"url_signature"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := s.jsonRequest(ctx, http.MethodGet, s.config.APIURL+"/file-upload", token, "", nil, &upload); err != nil {
		return "", fmt.Errorf("request Pingen upload: %w", err)
	}
	if upload.Data.Attributes.URL == "" || upload.Data.Attributes.Signature == "" {
		return "", &Error{operation: "file-upload response", temporary: true}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, upload.Data.Attributes.URL, bytes.NewReader(pdf))
	if err != nil {
		return "", &Error{operation: "PDF upload request"}
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := s.client.Do(req)
	if err != nil {
		return "", &Error{operation: "PDF upload", temporary: true}
	}
	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", statusError("PDF upload", resp.StatusCode, "")
	}

	payload := map[string]any{"data": map[string]any{
		"type": "letters",
		"attributes": map[string]any{
			"file_original_name": filename,
			"file_url":           upload.Data.Attributes.URL,
			"file_url_signature": upload.Data.Attributes.Signature,
			"address_position":   addressPosition,
			"auto_send":          true,
			"delivery_product":   product,
			"print_mode":         printMode,
			"print_spectrum":     printSpectrum,
		},
	}}
	var letter struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	endpoint := fmt.Sprintf("%s/organisations/%s/deliveries/letters", s.config.APIURL, url.PathEscape(s.config.OrganisationID))
	if err := s.jsonRequest(ctx, http.MethodPost, endpoint, token, idempotencyKey, payload, &letter); err != nil {
		return "", fmt.Errorf("create Pingen letter: %w", err)
	}
	if letter.Data.ID == "" {
		return "", &Error{operation: "letter response", temporary: true}
	}
	return letter.Data.ID, nil
}

func (s *Service) token(ctx context.Context) (string, error) {
	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {s.config.ClientID},
		"client_secret": {s.config.ClientSecret},
		"scope":         {s.config.Scope},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.config.IdentityURL+"/auth/access-tokens", strings.NewReader(form.Encode()))
	if err != nil {
		return "", &Error{operation: "token request"}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var response struct {
		AccessToken string `json:"access_token"`
	}
	if err := s.doJSON(req, &response); err != nil {
		return "", fmt.Errorf("obtain Pingen token: %w", err)
	}
	if response.AccessToken == "" {
		return "", &Error{operation: "token response", temporary: true}
	}
	return response.AccessToken, nil
}

func (s *Service) jsonRequest(ctx context.Context, method, endpoint, token, idempotencyKey string, body any, target any) error {
	var data bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&data).Encode(body); err != nil {
			return err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, &data)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.api+json")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	if idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", idempotencyKey)
	}
	return s.doJSON(req, target)
}

func (s *Service) doJSON(req *http.Request, target any) error {
	resp, err := s.client.Do(req)
	if err != nil {
		return &Error{operation: req.Method + " request", temporary: true}
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var response struct {
			Errors []struct {
				Code string `json:"code"`
			} `json:"errors"`
		}
		_ = json.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&response)
		code := ""
		if len(response.Errors) != 0 {
			code = response.Errors[0].Code
		}
		return statusError(req.Method+" request", resp.StatusCode, code)
	}
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return &Error{operation: req.Method + " response", temporary: true}
	}
	return nil
}

func statusError(operation string, status int, code string) *Error {
	code = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' || r == '-' {
			return r
		}
		return -1
	}, code)
	if len(code) > 64 {
		code = code[:64]
	}
	return &Error{operation: operation, status: status, code: code, temporary: status == 408 || status == 409 || status == 424 || status == 425 || status == 429 || status >= 500}
}
