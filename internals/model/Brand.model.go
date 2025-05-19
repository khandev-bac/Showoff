package model

import (
	"time"

	"github.com/google/uuid"
)

type Brand struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name       string    `gorm:"not null"`
	Logo       string    `gorm:"not null"`
	WebsiteURL string    `gorm:"not null"`
	Mission    string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
