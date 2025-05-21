package pkg

import (
	"context"
	"errors"
)

type ContextKey string

const USERID ContextKey = "userID"

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(USERID).(string)
	if !ok || userID == "" {
		return "", errors.New("user ID not found in context")
	}
	return userID, nil
}
