package models

import (
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Title     string    `gorm:"not null"`
	StartTime time.Time `gorm:"not null;index"`
	EndTime   time.Time `gorm:"not null;index"`
	CreatedAt time.Time
}
