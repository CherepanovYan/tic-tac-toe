package handler

import (
	"context"
	"tic-tac-toe/internal/domain/model"
)

type contextKey string

const userContextKey contextKey = "user"

func UserFromContext(ctx context.Context) (*model.User, bool) {
	user, ok := ctx.Value(userContextKey).(*model.User)
	return user, ok
}
