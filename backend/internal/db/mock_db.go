package database

import (
	"testing"

	"github.com/beejay1293/schedule-management/backend/internal/models"
	"gorm.io/gorm"

	"gorm.io/driver/sqlite"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// Auto-migrate tables
	if err := db.AutoMigrate(&models.Appointment{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	// Create unique index to simulate Postgres constraint
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS unique_appointment_time 
		ON appointments(date);
	`).Error; err != nil {
		t.Fatalf("failed to create unique index: %v", err)
	}

	return db
}
