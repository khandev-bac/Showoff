package model

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	BrandID     uuid.UUID `gorm:"type:uuid;not null"`
	Brand       Brand     `gorm:"foreignKey:BrandID"`
	ProductName string    `gorm:"not null"`
	ProductLink string    `gorm:"not null"`
	Caption     string
	Images      []PostImage `gorm:"foreignKey:PostID"`
	CreatedAt   time.Time   `gorm:"autoCreateTime"`
}
type PostImage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	PostID    uuid.UUID `gorm:"type:uuid;not null"`
	Post      Post      `gorm:"foreignKey:PostID"`
	ImageURL  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
