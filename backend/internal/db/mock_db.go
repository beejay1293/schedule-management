package database

import (
	"testing"

	"github.com/beejay1293/schedule-management/backend/internal/models"
	"gorm.io/gorm"

	"gorm.io/driver/sqlite"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	// Use shared in-memory database
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// Run migrations
	if err := db.AutoMigrate(&models.Appointment{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	// Add unique index for conflict test
	if err := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS unique_appointment_time ON appointments(date);`).Error; err != nil {
		t.Fatalf("failed to create unique index: %v", err)
	}

	// Force all goroutines to reuse the same single connection
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	return db
}
