package storage

import (
	"database/sql"
	"kai-sec/internal/logger"

	// SQLite driver for testing
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

var TestDB *sql.DB
var log = logger.GetLogger()

// InitTestDB initializes an in-memory SQLite DB for unit testing
func InitTestDB() {
	log.Info("Initializing in-memory test database...")

	var err error
	TestDB, err = sql.Open("sqlite3", ":memory:") // Create in-memory DB
	if err != nil {
		log.Fatal("Failed to initialize test DB", zap.Error(err))
	}

	// Define table creation queries
	createTables := []string{
		`CREATE TABLE scans (
			scan_id TEXT PRIMARY KEY,
			timestamp TIMESTAMP,
			scan_status TEXT,
			resource_type TEXT,
			resource_name TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_attempted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE vulnerabilities (
			id TEXT PRIMARY KEY,
			scan_id TEXT,
			severity TEXT,
			cvss REAL,
			status TEXT,
			package_name TEXT,
			current_version TEXT,
			fixed_version TEXT,
			description TEXT,
			published_date TIMESTAMP,
			link TEXT,
			risk_factors TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(scan_id) REFERENCES scans(scan_id) ON DELETE CASCADE
		);`,
		`CREATE TABLE scan_summary (
			scan_id TEXT PRIMARY KEY,
			total_vulnerabilities INTEGER,
			severity_counts TEXT,
			fixable_count INTEGER,
			compliant BOOLEAN,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(scan_id) REFERENCES scans(scan_id) ON DELETE CASCADE
		);`,
		`CREATE TABLE scan_metadata (
			scan_id TEXT PRIMARY KEY,
			scanner_version TEXT,
			policies_version TEXT,
			scanning_rules TEXT,
			excluded_paths TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(scan_id) REFERENCES scans(scan_id) ON DELETE CASCADE
		);`,
	}

	// Execute all table creation queries
	for _, query := range createTables {
		_, err = TestDB.Exec(query)
		if err != nil {
			log.Fatal("Failed to create test tables", zap.Error(err))
		}
	}

	log.Info("Test database initialized with required tables.")
}
