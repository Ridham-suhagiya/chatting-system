package model

import (
	"time"

	"github.com/google/uuid"
)

type ChatLinks struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserId    uuid.UUID `gorm:"type:uuid;not null;foreignKey:users(id)"`
	LinkCode  string    `gorm:"unique;not null"`
	ExpiryAt  time.Time `gorm:"type:timestamp;"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	NumUsers  int       `gorm:"not null;type:integer;default:0"`
}
