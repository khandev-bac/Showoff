package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"unique;not null"`
	Password     *string
	ProfilePic   string
	RefreshToken string
	Bio          string
	GoogleID     *string
	IsOauthUser  bool
	Likes        []Like     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	WishLists    []WishList `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	CreatedAt    time.Time
}
