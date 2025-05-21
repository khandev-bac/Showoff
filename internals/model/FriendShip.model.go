package model

import (
	"time"

	"github.com/google/uuid"
)

type FriendShip struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	User1ID   uuid.UUID `gorm:"type:uuid;not null;index"`
	User1     User      `gorm:"foreignKey:User1ID;constraint:OnDelete:CASCADE"`
	User2ID   uuid.UUID `gorm:"type:uuid;not null;index"`
	User2     User      `gorm:"foreignKey:User2ID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time
}
