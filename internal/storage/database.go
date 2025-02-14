package storage

import (
	"database/sql"
	"fmt"
	"kai-sec/internal/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

var DB *sql.DB
var l = logger.GetLogger()

func ConnectDB() (*sql.DB, error) {
	l.Info("Initializing database connection........")
	//TODO: Make it env variable
	dsn := "postgres://user:password@localhost:5432/vulnerabilities_db?sslmode=disable&search_path=public"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		l.Error("Failed to open database connection", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	err = db.Ping()
	if err != nil {
		l.Error("Database ping failed", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	l.Info("Database connected successfully!")

	// Run migrations
	if err := RunMigrations(db); err != nil {
		l.Error("Database migration failed", zap.Error(err))
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	l.Info("Database connection established and migrations applied successfully!")
	return db, nil
}

func RunMigrations(db *sql.DB) error {
	l.Info("Applying database migrations...")
	err := goose.Up(db, "migrations", goose.WithAllowMissing())
	if err != nil {
		l.Error("Goose migration failed", zap.Error(err))
		return fmt.Errorf("goose migration failed: %w", err)
	}
	l.Info("All migrations applied successfully!")
	return nil
}
