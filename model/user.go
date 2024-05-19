package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id         uuid.UUID `gorm:"primaryKey; default: uuid_generate_v4()"`
	Username   string    `gorm:"unique;not null"`
	Email      string    `gorm:"unique;not null"`
	Password   string    `gorm:"not null; unique"`
	Created_at time.Time `gorm:"not null; default:CURRENT_TIMESTAMP"`
}
