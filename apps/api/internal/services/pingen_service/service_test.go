package pingen_service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubmit(t *testing.T) {
	var server *httptest.Server
	tokenCalls, letterCalls := 0, 0
	reject := false
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/access-tokens":
			tokenCalls++
			if tokenCalls == 1 {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			require.NoError(t, r.ParseForm())
			assert.Equal(t, url.Values{
				"client_id": {"client"}, "client_secret": {"secret"},
				"grant_type": {"client_credentials"}, "scope": {"letter"},
			}, r.Form)
			_, _ = io.WriteString(w, `{"access_token":"token"}`)
		case "/file-upload":
			assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))
			_, _ = io.WriteString(w, `{"data":{"attributes":{"url":"`+server.URL+`/upload","url_signature":"signature"}}}`)
		case "/upload":
			data, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			assert.Equal(t, []byte("pdf"), data)
		case "/organisations/org/deliveries/letters":
			letterCalls++
			assert.Equal(t, "delivery-id", r.Header.Get("Idempotency-Key"))
			var payload struct {
				Data struct {
					Attributes map[string]any `json:"attributes"`
				} `json:"data"`
			}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			assert.Equal(t, true, payload.Data.Attributes["auto_send"])
			assert.Equal(t, "fast", payload.Data.Attributes["delivery_product"])
			assert.Equal(t, "left", payload.Data.Attributes["address_position"])
			if reject {
				w.WriteHeader(http.StatusUnprocessableEntity)
				_, _ = io.WriteString(w, `{"errors":[{"code":"address_invalid","detail":"recipient secret"}]}`)
				return
			}
			if letterCalls == 1 {
				connection, _, err := w.(http.Hijacker).Hijack()
				require.NoError(t, err)
				_ = connection.Close()
				return
			}
			_, _ = io.WriteString(w, `{"data":{"id":"pingen-id"}}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	service := New(config.PingenConfig{
		ClientID: "client", ClientSecret: "secret", OrganisationID: "org",
		IdentityURL: server.URL, APIURL: server.URL, Scope: "letter",
	})
	_, err := service.Submit(context.Background(), []byte("pdf"), "letter.pdf", "delivery-id", "fast", "left", "simplex", "color")
	var temporary interface{ Temporary() bool }
	require.ErrorAs(t, err, &temporary)
	assert.True(t, temporary.Temporary())

	id, err := service.Submit(context.Background(), []byte("pdf"), "letter.pdf", "delivery-id", "fast", "left", "simplex", "color")
	require.NoError(t, err)
	assert.Equal(t, "pingen-id", id)
	assert.Equal(t, 2, letterCalls)

	reject = true
	_, err = service.Submit(context.Background(), []byte("pdf"), "letter.pdf", "delivery-id", "fast", "left", "simplex", "color")
	var pingenError *Error
	require.ErrorAs(t, err, &pingenError)
	assert.False(t, pingenError.Temporary())
	assert.Equal(t, "pingen_status_422_address_invalid", pingenError.SafeMessage())
	assert.NotContains(t, pingenError.SafeMessage(), "recipient secret")
}
