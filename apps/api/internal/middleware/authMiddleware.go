package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/bored-bots/bored-template/internal/context"
	"github.com/bored-bots/bored-template/internal/errorcodes"
	"github.com/danielgtaylor/huma/v2"
)

const RoleAdmin = "admin"

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
		Role  string `json:"role"`
	} `json:"user"`
}

type UserContext struct {
	ID    string
	Email string
	Name  string
	Role  string
}

func NewAuthMiddleware(api huma.API, authConfig *config.AuthConfig) func(ctx huma.Context, next func(ctx2 huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		schemes := ctx.Operation().Security

		if len(schemes) == 0 {
			next(ctx)
			return
		}

		requiresSessionCookie := false
		requiresAPIKey := false
		for _, schemeMap := range schemes {
			if _, exists := schemeMap["session-cookie"]; exists {
				requiresSessionCookie = true
			}
			if _, exists := schemeMap["apiKeyAuth"]; exists {
				requiresAPIKey = true
			}
		}

		if !requiresSessionCookie && !requiresAPIKey {
			next(ctx)
			return
		}

		var sessionCookieValue string
		var apiKeyValue string
		var err error

		if requiresSessionCookie {
			sessionCookie, err := huma.ReadCookie(ctx, authConfig.CookieName)
			if err == nil {
				sessionCookieValue = sessionCookie.Value
			}
		}

		if requiresAPIKey {
			apiKeyValue = ctx.Header("X-API-Key")
		}

		if sessionCookieValue == "" && apiKeyValue == "" {
			log.Println("No valid authentication credentials found")
			huma.WriteErr(api, ctx, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		req, err := http.NewRequest("GET", authConfig.AuthServerURL+"/api/auth/get-session", nil)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusInternalServerError, errorcodes.AuthSessionRequestBuild)
			return
		}

		if sessionCookieValue != "" {
			req.Header.Set("Cookie", authConfig.CookieName+"="+sessionCookieValue)
		}
		if apiKeyValue != "" {
			req.Header.Set("X-API-Key", apiKeyValue)
		}

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
			Role:  sessionResp.User.Role,
		}

		ctx = context.WithUserID(ctx, userCtx.ID)

		for _, schemeMap := range schemes {
			if requiredRoles, exists := schemeMap["session-cookie"]; exists {
				if len(requiredRoles) > 0 {

					if userCtx.Role == RoleAdmin {
						break
					}

					hasRequiredRole := false
					for _, requiredRole := range requiredRoles {
						if userCtx.Role == requiredRole {
							hasRequiredRole = true
							break
						}
					}
					if !hasRequiredRole {
						log.Printf("User %s does not have required role. User role: %s, Required roles: %v", userCtx.ID, userCtx.Role, requiredRoles)
						huma.WriteErr(api, ctx, http.StatusForbidden, "Insufficient permissions")
						return
					}
				}
			}
		}

		next(ctx)
	}
}
