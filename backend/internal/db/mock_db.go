package database

import (
	"fmt"
	"log"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	dsnBase := "host=localhost user=postgres password=postgres port=5432 sslmode=disable dbname=postgres"
	baseDB, err := gorm.Open(postgres.Open(dsnBase), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to Postgres: %v", err)
	}

	testDBName := "schedule_db_test"

	// Check if database exists
	var exists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s');", testDBName)
	if err := baseDB.Raw(checkQuery).Scan(&exists).Error; err != nil {
		t.Fatalf("failed to check if test database exists: %v", err)
	}

	// Create database only if it doesn't exist
	if !exists {
		if err := baseDB.Exec(fmt.Sprintf("CREATE DATABASE %s;", testDBName)).Error; err != nil {
			t.Fatalf("failed to create test database: %v", err)
		}
		log.Printf("Created test database: %s", testDBName)
	} else {
		log.Printf("Test database %s already exists", testDBName)
	}

	// Connect to the test database
	dsnTest := fmt.Sprintf("host=localhost user=postgres password=postgres port=5432 sslmode=disable dbname=%s", testDBName)
	db, err := gorm.Open(postgres.Open(dsnTest), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Run migrations
	if err := Migrate(db); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}

	log.Printf("✅ Test database %s is ready and migrated", testDBName)

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()

		dropQuery := fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE);", testDBName)
		if err := baseDB.Exec(dropQuery).Error; err != nil {
			t.Logf("⚠️ Failed to drop test database %s: %v", testDBName, err)
		} else {
			t.Logf("🗑️ Dropped test database %s after test cleanup.", testDBName)
		}
	})

	return db
}
