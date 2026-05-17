package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/bored-bots/bored-template/internal/errorcodes"
	"github.com/danielgtaylor/huma/v2"
)

type contextKey string

var (
	ContextKeyUserID = contextKey("user-id")
)

type SessionResponse struct {
	Session struct {
		UserID string `json:"userId"`
		Email  string `json:"email"`
		Name   string `json:"name"`
	} `json:"session"`
	User struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"user"`
}

type UserContext struct {
	ID    string
	Email string
	Name  string
}

func NewAuthMiddleware(api huma.API, authConfig *config.AuthConfig) func(ctx huma.Context, next func(ctx2 huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		schemes := ctx.Operation().Security

		if len(schemes) == 0 {
			next(ctx)
			return
		}

		requiresSessionCookie := false
		for _, schemeMap := range schemes {
			if _, exists := schemeMap["session-cookie"]; exists {
				requiresSessionCookie = true
				break
			}
		}

		if requiresSessionCookie {
			sessionCookie, err := huma.ReadCookie(ctx, authConfig.CookieName)
			if err != nil {
				log.Println("Session cookie not found")
				if err != nil {
					huma.WriteErr(api, ctx, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
					return
				}
			}

			req, err := http.NewRequest("GET", authConfig.AuthServerURL+"/api/auth/get-session", nil)
			if err != nil {
				huma.WriteErr(api, ctx, http.StatusInternalServerError, errorcodes.AuthSessionRequestBuild)
				return
			}

			req.Header.Set("Cookie", authConfig.CookieName+"="+sessionCookie.Value)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				huma.WriteErr(api, ctx, http.StatusInternalServerError, errorcodes.AuthSessionValidateFail)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				huma.WriteErr(api, ctx, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}

			var sessionResp SessionResponse
			if err := json.NewDecoder(resp.Body).Decode(&sessionResp); err != nil {
				huma.WriteErr(api, ctx, http.StatusInternalServerError, errorcodes.AuthSessionParseFail)
				return
			}
			if sessionResp.Session.UserID == "" {
				huma.WriteErr(api, ctx, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}

			userCtx := &UserContext{
				ID:    sessionResp.User.ID,
				Email: sessionResp.User.Email,
				Name:  sessionResp.User.Name,
			}

			ctx = huma.WithValue(ctx, ContextKeyUserID, userCtx.ID)
		}

		next(ctx)
	}
}
