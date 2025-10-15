package database

import "gorm.io/gorm"

type DB interface {
	GetDB() *gorm.DB
	Close() error
}

type PostgresDB struct {
	db *gorm.DB
}

func NewPostgresDB() (DB, error) {
	gormDB := ProvidePostgres()
	return &PostgresDB{db: gormDB}, nil
}

func (p *PostgresDB) GetDB() *gorm.DB {
	return p.db
}

func (p *PostgresDB) Close() error {
	return Close(p.db)
}
