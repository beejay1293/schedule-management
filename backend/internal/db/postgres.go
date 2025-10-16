package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ProvidePostgres returns a GORM postgres database connection
// and optionally runs migrations if RUN_DB_MIGRATIONS=true
func ProvidePostgres() *gorm.DB {
	postgresDbUrl := os.Getenv("POSTGRES_DB_URL")

	// Connect to the postgres database
	gormdb, err := Connect(postgres.Open(postgresDbUrl))
	if err != nil {
		log.Fatalf("postgres database connection error: %v", err)
	}

	fmt.Println("run migration", os.Getenv("RUN_DB_MIGRATIONS"))

	// // Optionally run migrations
	if os.Getenv("RUN_DB_MIGRATIONS") == "true" {
		if err := Migrate(gormdb); err != nil {
			log.Fatalf("failed to run database migrations: %v", err)
		}
	}

	return gormdb
}

// Connect provides functionality for connecting to the database
func Connect(dialector gorm.Dialector) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Close provides functionality for closing the database connection
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
