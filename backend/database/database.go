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

//initializes the database connection and runs migrations
func ConnectDB(cfg *config.Config) {
	var err error

	//retries database connection
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

	//run migrations
	if err := runMigrations(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	fmt.Println("Database migrations completed successfully!")
}

func runMigrations() error {
	var exists bool
	err := DB.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'users')",
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("error checking for existing tables: %v", err)
	}

	if exists {
		fmt.Println("Tables already exist, skipping migrations.")
		return nil
	}

	//run migrations only if the tables do not exist
	migrationsPath := "/app/database/migrations.sql"
	migrations, err := os.ReadFile(migrationsPath)
	if err != nil {
		return fmt.Errorf("error reading migrations file: %v", err)
	}

	_, err = DB.Exec(context.Background(), string(migrations))
	if err != nil {
		return fmt.Errorf("error executing migrations: %v", err)
	}

	fmt.Println("Database migrations completed successfully!")
	return nil
}
