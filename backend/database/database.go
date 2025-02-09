package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"storgage/config"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

// ConnectDB initializes the database connection and runs migrations
func ConnectDB(cfg *config.Config) {
	var err error
	DB, err = pgx.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	fmt.Println("✅ Database connected successfully!")

	//run migrations
	if err := runMigrations(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	fmt.Println("✅ Database migrations completed successfully!")
}

func runMigrations() error {
	wd, _ := os.Getwd()
	fmt.Println("📂 Current working directory:", wd)

	_, err := os.Stat("database/migrations.sql")
	if err != nil {
		fmt.Println("❌ migrations.sql NOT found in database/ directory")
		return fmt.Errorf("error finding migrations file: %v", err)
	} else {
		fmt.Println("✅ migrations.sql FOUND in database/ directory")
	}

	migrations, err := os.ReadFile("database/migrations.sql")
	if err != nil {
		return fmt.Errorf("error reading migrations file: %v", err)
	}

	_, err = DB.Exec(context.Background(), string(migrations))
	if err != nil {
		return fmt.Errorf("error executing migrations: %v", err)
	}

	fmt.Println("✅ Database migrations completed successfully!")
	return nil
}
