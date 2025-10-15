package main

import (
	"log"

	"github.com/beejay1293/schedule-management/backend/internal/app"
	database "github.com/beejay1293/schedule-management/backend/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// Connect to Postgres
	db, err := database.NewPostgresDB()
	if err != nil {

	}
	defer db.Close()

	provideDB := db.GetDB()

	// Auto-migrate appointment table
	if err := database.Migrate(provideDB); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create and initialize the app
	a, err := app.New(provideDB)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	// Run the gRPC server
	a.Run()
}
