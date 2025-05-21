package service

import (
	"context"
	"errors"
	"exceapp/internals/model"
	"exceapp/internals/repo"

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

func (s *UserService) Signup(ctx context.Context, user *model.User) (*model.User, error) {
	exist, err := s.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		return nil, errors.New("internal error while checking user")
	}
	if exist != nil {
		return nil, errors.New("user already exists")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	newUser := &model.User{
		ID:       uuid.New(),
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hashedPass),
	}

	return s.repo.Signup(ctx, newUser)
}

func (s *UserService) Login(ctx context.Context, email, password string) (*model.User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *model.User) error {
	return s.repo.UpdateUser(ctx, user)
}
func (s *UserService) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	return s.repo.UpdateRefreshToken(ctx, userID, refreshToken)
}
func (s *UserService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.FindByEmail(ctx, email)
}
func (s *UserService) FindById(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	return s.repo.FindById(ctx, userID)
}
func (s *UserService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteUser(ctx, userID)
}
