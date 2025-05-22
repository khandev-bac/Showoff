package service

import (
	"context"
	"exceapp/internals/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SwipeService struct {
	db *gorm.DB
}

func NewSwipeService(db *gorm.DB) *SwipeService {
	return &SwipeService{
		db: db,
	}
}
func (sw *SwipeService) GetUnswippedUsers(ctx context.Context, currentId uuid.UUID) (*model.User, error) {
	var swipedIDs []uuid.UUID
	sw.db.WithContext(ctx).Model(&model.Swipe{}).Where("swiper_id = ?", currentId).Pluck("swiped_id", swipedIDs)
	var user model.User
	err := sw.db.WithContext(ctx).Model(&model.User{}).Where("id != ?", currentId).Where("id NOT IN ?", swipedIDs).Order("RANDOM()").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (sw *SwipeService) SaveSwipes(ctx context.Context, SwiperID, SwipedID uuid.UUID) error {
	swipe := &model.Swipe{
		ID:        uuid.New(),
		SwiperID:  SwiperID,
		SwipedID:  SwipedID,
		CreatedAt: time.Now(),
	}
	return sw.db.WithContext(ctx).Create(swipe).Error
}

func (sw *SwipeService) GetSwippedHistory(ctx context.Context, userID uuid.UUID) ([]model.User, error) {
	var swipes []model.Swipe
	if err := sw.db.WithContext(ctx).Preload("Swiped").Where("swiper_id = ?", userID).Find(&swipes).Error; err != nil {
		return nil, err
	}
	users := make([]model.User, 0, len(swipes))
	for _, swipe := range swipes {
		users = append(users, swipe.Swiped)
	}
	return users, nil
}
