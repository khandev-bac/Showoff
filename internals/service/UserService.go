package service

import (
	"context"
	"errors"
	"exceapp/internals/model"
	"exceapp/internals/repo"
	"exceapp/pkg/jwt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repo.UserRepo
}

func NewUserService(repo *repo.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}
func (s *UserService) Signup(ctx context.Context, user *model.User) error {
	exist, _ := s.repo.GetUserByEmail(ctx, user.Email)
	if exist != nil {
		return errors.New("user already exists")
	}
	hasedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := &model.User{
		ID:        uuid.New(),
		Name:      user.Name,
		Email:     user.Email,
		Password:  string(hasedPass),
		CreatedAt: time.Now(),
	}
	return s.repo.CreateUser(ctx, newUser)
}
func (s *UserService) Login(ctx context.Context, email, password string) (*model.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}
	token, err := jwt.GenerateJWTToken(user.ID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.UpdateRefreshToken(ctx, user.ID, token.RefreshToken); err != nil {
		return nil, err
	}
	return user, nil

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
