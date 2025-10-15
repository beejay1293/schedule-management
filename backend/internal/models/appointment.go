package models

import "time"

type Appointment struct {
	ID        string    `gorm:"primaryKey"`
	Title     string    `gorm:"not null"`
	Date      time.Time `gorm:"not null;index"`
	CreatedAt time.Time
}
