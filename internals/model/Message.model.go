package model

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	FromUserID uuid.UUID `gorm:"type:uuid;not null"`
	FromUser   User      `gorm:"foreignKey:FromUserID;constraint:OnDelete:CASCADE"`

	ToUserID uuid.UUID `gorm:"type:uuid;not null"`
	ToUser   User      `gorm:"foreignKey:ToUserID;constraint:OnDelete:CASCADE"`

	Text      string `gorm:"type:text"`
	CreatedAt time.Time
}
