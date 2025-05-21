package model

import (
	"time"

	"github.com/google/uuid"
)

type Swipe struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	FromUserID uuid.UUID `gorm:"type:uuid;not null"`
	FromUser   User      `gorm:"foreignKey:FromUserID;constraint:OnDelete:CASCADE"`
	ToUserID   uuid.UUID `gorm:"type:uuid;not null"`
	ToUser     User      `gorm:"foreignKey:FromUserID;constraint:OnDelete:CASCADE"`
	Direction  string    `gorm:"not null;index"`
	Emoji      string
	CreatedAt  time.Time
}
