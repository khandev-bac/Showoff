package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	Password     string
	ProfilePic   string
	RefreshToken string
	GoogleID     string `gorm:"index"`
	IsOauthUser  bool
	CreatedAt    time.Time
}
