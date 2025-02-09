package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"storgage/config"
	"time"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

// ConnectDB initializes the database connection and runs migrations
func ConnectDB(cfg *config.Config) {
	var err error

	// Retry database connection
	for i := 0; i < 10; i++ {
		DB, err = pgx.Connect(context.Background(), cfg.DatabaseURL)
		if err == nil {
			fmt.Println("Database connected successfully!")
			break
		}
		fmt.Println("Waiting for database to be ready... Retrying in 5s")
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	// Run migrations
	if err := runMigrations(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	fmt.Println("Database migrations completed successfully!")
}

func runMigrations() error {
	migrationsPath := "/app/database/migrations.sql"

	// Print working directory (debugging)
	wd, _ := os.Getwd()
	fmt.Println("Current working directory:", wd)

	// Check if migrations.sql exists
	_, err := os.Stat(migrationsPath)
	if err != nil {
		fmt.Println("migrations.sql NOT found at:", migrationsPath)
		return fmt.Errorf("error finding migrations file: %v", err)
	} else {
		fmt.Println("migrations.sql FOUND at:", migrationsPath)
	}

	// Read migrations file
	migrations, err := os.ReadFile(migrationsPath)
	if err != nil {
		return fmt.Errorf("error reading migrations file: %v", err)
	}

	// Execute migrations
	_, err = DB.Exec(context.Background(), string(migrations))
	if err != nil {
		return fmt.Errorf("error executing migrations: %v", err)
	}

	fmt.Println("Database migrations completed successfully!")
	return nil
}
