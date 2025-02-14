package daos

import (
	"encoding/json"
	"fmt"
	"kai-sec/internal/daos/models"
	"kai-sec/internal/logger"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
)

func TestUpsertScanMetadata(t *testing.T) {
	l := logger.GetLogger()

	db, mock, err := sqlmock.New()
	if err != nil {
		l.Fatal("Failed to create mock DB", zap.Error(err))
	}
	defer db.Close()

	// start a transaction expectation
	mock.ExpectBegin()

	// Sample test data
	metadata := models.ScanMetadataDAO{
		ScannerVersion:  "30.1.51",
		PoliciesVersion: "2025.1.29",
		ScanningRules:   []string{"vulnerability", "compliance", "malware"},
		ExcludedPaths:   []string{"/tmp", "/var/log"},
	}

	l.Info("starting upsertscanmetadata test", zap.String("scanner_version", metadata.ScannerVersion))

	scanningRulesJSON, _ := json.Marshal(metadata.ScanningRules)
	excludedPathsJSON, _ := json.Marshal(metadata.ExcludedPaths)

	// Expect INSERT INTO scan_metadata
	mock.ExpectExec(`INSERT INTO scan_metadata`).
		WithArgs("scan_001", metadata.ScannerVersion, metadata.PoliciesVersion, string(scanningRulesJSON), string(excludedPathsJSON)).
		WillReturnResult(sqlmock.NewResult(0, 1)) // No LastInsertId in PostgreSQL

	// Expect transaction commit
	mock.ExpectCommit()

	// Start a transaction in mock DB
	tx, err := db.Begin()
	if err != nil {
		l.Fatal("failed to start transaction", zap.Error(err))
	}

	// Run function with transaction
	err = UpsertScanMetadata(tx, "scan_001", metadata)
	if err != nil {
		l.Error("failed to insert scan metadata", zap.Error(err))
		t.Fatalf("failed to insert scan metadata: %v", err)
	}

	// Commit transaction explicitly
	err = tx.Commit()
	if err != nil {
		t.Fatalf("failed to commit transaction: %v", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("transaction expectations not met: %v", err)
	}
}
func TestUpsertScanMetadata_Error(t *testing.T) {
	l := logger.GetLogger()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Expect transaction begin
	mock.ExpectBegin()

	metadata := models.ScanMetadataDAO{
		ScannerVersion:  "1.2.3",
		PoliciesVersion: "2025.2.10",
		ScanningRules:   []string{"test-rule"},
		ExcludedPaths:   []string{"/test/path"},
	}

	// Expect INSERT query but return an error
	mock.ExpectExec("INSERT INTO scan_metadata").
		WithArgs("scan_test", metadata.ScannerVersion, metadata.PoliciesVersion, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(fmt.Errorf("mock insert failure"))

	// Expect transaction rollback
	mock.ExpectRollback()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to start transaction: %v", err)
	}

	l.Info("starting error case for upsertscanmetadata test")

	err = UpsertScanMetadata(tx, "scan_test", metadata)
	if err == nil {
		t.Errorf("expected error, got nil")
	} else {
		l.Info("expected failure triggered successfully", zap.Error(err))
	}

	// Explicitly rollback transaction
	err = tx.Rollback()
	if err != nil {
		t.Fatalf("failed to rollback transaction: %v", err)
	}

	// Ensure expectations met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("transaction expectations not met: %v", err)
	}
}