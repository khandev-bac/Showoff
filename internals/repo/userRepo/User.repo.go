package userrepo

import (
	"context"
	"exceapp/internals/model"

	"github.com/google/uuid"
)

type UserRepoInterface interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, name string, profilepic string) error
	DeleteAccount(ctx context.Context, userID uuid.UUID) error
	UpdateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) error
}
