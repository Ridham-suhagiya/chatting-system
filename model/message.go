package model

import (
	"time"

	"github.com/google/uuid"
)

type Messages struct {
	MessageId      uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	MessageContent []byte    `gorm:"type:jsonb"`
	LinkId         uuid.UUID `gorm:"type:uuid;not null;foreignKey:ChatLinks(id)"`
	Timestamp      time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
