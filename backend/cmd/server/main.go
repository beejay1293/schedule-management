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
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	provideDB := db.GetDB()

	// initialize the app
	a, err := app.New(provideDB)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	// Run server
	a.Run()
}
