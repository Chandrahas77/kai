package storage

import (
	"database/sql"
	"fmt"
	"kai-sec/internal/config"
	"kai-sec/internal/logger"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

var DB *sql.DB
var l = logger.GetLogger()

// InitDB initializes a database connection and returns a *sql.DB instance
func InitDB(cfg *config.DBConfig) (*sql.DB, error) {
	l.Info("Initializing database connection........")

	dsn := cfg.GetDSN()

	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			l.Error("Failed to open database connection", zap.Error(err))
		} else {
			err = db.Ping()
			if err == nil {
				break
			}
			l.Error("Database ping failed", zap.Error(err))
		}
		time.Sleep(2 * time.Second)
	}
	DB = db
	l.Info("Database connected successfully!", zap.Any("db", DB))

	// Run migrations
	err = RunMigrations(db)
	if err != nil {
		l.Error("Database migration failed", zap.Error(err))
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

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