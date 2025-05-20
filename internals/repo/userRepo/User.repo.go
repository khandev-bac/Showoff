package userrepo

import (
	"context"
	"exceapp/internals/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, name string, profilepic string) error
	DeleteAccount(ctx context.Context, userID uuid.UUID) error
	UpdateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) error
	GetByGoogleID(ctx context.Context, googleID string) (*model.User, error)
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}
func (r *UserRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	return &user, err
}
func (r *UserRepo) GetByGoogleID(ctx context.Context, googleID uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("google_id = ?", googleID).First(&user).Error
	return &user, err
}
func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}
func (r *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) UpdateUser(ctx context.Context, userID uuid.UUID, name string, profilepic string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"name":        name,
		"profile_pic": profilepic,
	}).Error
}
func (r *UserRepo) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", userID).Delete(&model.User{}).Error
}
func (r *UserRepo) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Update("refresh_token", refreshToken).Error
}
