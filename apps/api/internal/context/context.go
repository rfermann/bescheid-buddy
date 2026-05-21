package context

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type contextKey string

var (
	UserID = contextKey("user-id")
)

func WithUserID(ctx huma.Context, userID string) huma.Context {
	return huma.WithValue(ctx, UserID, userID)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, err := uuid.Parse(ctx.Value(UserID).(string))
	if err != nil {
		return [16]byte{}, false
	}

	return userID, true
}
