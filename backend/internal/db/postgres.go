package database

import (
	"os"
	"time"

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
		// zerolog.Logger.Err(err).Msg("postgres database connection error")
	}

	// Optionally run migrations
	if os.Getenv("RUN_DB_MIGRATIONS") == "true" {
		if err := Migrate(gormdb); err != nil {
			// zerolog.Fatal().Err(err).Msg("failed to run database migrations")
		}
	}

	// Configure connection pool
	sqlDB, err := gormdb.DB()
	if err != nil {
		// zerolog.Fatal().Err(err).Msg("failed to get sql.DB instance")
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

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
