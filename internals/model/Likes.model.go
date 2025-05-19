package model

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	User      User      `gorm:"foreignKey:UserID"`
	PostID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Post      Post      `gorm:"foreignKey:PostID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
