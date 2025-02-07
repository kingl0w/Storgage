package database

import (
	"context"
	"fmt"
	"log"

	"storgage/config"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

// ConnectDB initializes the database connection
func ConnectDB(cfg *config.Config) {
	var err error
	DB, err = pgx.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	fmt.Println("✅ Database connected successfully!")
}
