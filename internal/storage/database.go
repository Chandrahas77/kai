package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

var DB *sql.DB

func ConnectDB() (*sql.DB, error) {
	//TODO: Make it env variable
	dsn := "postgres://user:password@localhost:5432/vulnerabilities_db?sslmode=disable&search_path=public"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("Goose migration failed: %v", err)
	}

	fmt.Println("Connected to PostgreSQL and migrations applied!")
	return db, nil
}

func RunMigrations(db *sql.DB) error {
	if err := goose.Up(db, "migrations", goose.WithAllowMissing()); err != nil {
		return fmt.Errorf("goose migration failed: %w", err)
	}
	fmt.Println("Goose migrations applied successfully!")
	return nil
}
