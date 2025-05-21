package repo

import (
	"context"
	"exceapp/internals/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}
func (r *UserRepo) Signup(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.db.WithContext(ctx).Create(user).Error
	return user, err
}
func (r *UserRepo) UpdateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}
func (r *UserRepo) FindById(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not found = no user
		}
		return nil, err // Some other error
	}
	return &user, nil
}
func (r *UserRepo) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("refresh_token", refreshToken).Error
}
func (r *UserRepo) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? ", userID).Delete(&model.User{}).Error
}
func (r *UserRepo) FindAllUsers(ctx context.Context, limit, offset int, excludeUserID uuid.UUID) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).
		Where("id != ?", excludeUserID). // exclude current user
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
