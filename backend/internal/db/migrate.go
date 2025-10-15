package database

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

//go:embed migrations
var migrationFiles embed.FS

// Migrate runs all up migrations
func Migrate(db *gorm.DB) error {
	sourceDriver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return errors.Wrap(err, "failed to create migration source driver")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return errors.Wrap(err, "failed to get sql.DB from gorm.DB")
	}

	databaseDriver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "failed to create migration database driver")
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", databaseDriver)
	if err != nil {
		return errors.Wrap(err, "failed to create migration instance")
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "failed to run migrations")
	}

	return nil
}
