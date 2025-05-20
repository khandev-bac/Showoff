package service

import (
	"context"
	"exceapp/internals/model"
	"exceapp/internals/repo"

	"github.com/google/uuid"
)

type UserService struct {
	repo *repo.UserRepo
}

func NewUserService(repo *repo.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}
func (s *UserService) RegisterUser(ctx context.Context, user *model.User) error {
	return s.repo.CreateUser(ctx, user)
}
func (s *UserService) FindByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}
func (s *UserService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}
func (s *UserService) FindByGoogleID(ctx context.Context, userID string) (*model.User, error) {
	return s.repo.GetByGoogleID(ctx, userID)
}
func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, name, profilePic string) error {
	return s.repo.UpdateUser(ctx, userID, name, profilePic)
}
func (s *UserService) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteAccount(ctx, userID)
}

func (s *UserService) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	return s.repo.UpdateRefreshToken(ctx, userID, refreshToken)
}
