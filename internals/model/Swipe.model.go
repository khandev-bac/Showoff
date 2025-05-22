package model

import (
	"time"

	"github.com/google/uuid"
)

type Swipe struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	SwiperID uuid.UUID `gorm:"type:uuid;not null"`
	Swiper   User      `gorm:"foreignKey:SwiperID;references:ID;constraint:OnDelete:CASCADE"`

	SwipedID uuid.UUID `gorm:"type:uuid;not null"`
	Swiped   User      `gorm:"foreignKey:SwipedID;references:ID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time
}
