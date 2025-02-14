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

func TestUpsertScanSummary(t *testing.T) {
	l := logger.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		l.Fatal("failed to create mock DB", zap.Error(err))
	}
	defer db.Close()

	// Expect transaction start
	mock.ExpectBegin()

	// Sample test data
	summary := models.ScanSummaryDAO{
		TotalVulnerabilities: 3,
		SeverityCounts: models.SeverityCountsDAO{
			Critical: 1,
			High:     1,
			Medium:   1,
			Low:      0,
		},
		FixableCount: 2,
		Compliant:    false,
	}

	severityCountsJSON, _ := json.Marshal(summary.SeverityCounts)

	// Expect INSERT INTO scan_summary
	mock.ExpectExec("INSERT INTO scan_summary").
		WithArgs("scan_001", summary.TotalVulnerabilities, string(severityCountsJSON), summary.FixableCount, summary.Compliant).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect transaction commit
	mock.ExpectCommit()

	// Start a transaction in mock DB
	tx, err := db.Begin()
	if err != nil {
		l.Fatal("failed to start transaction", zap.Error(err))
	}

	// Run function
	err = UpsertScanSummary(tx, "scan_001", summary)
	if err != nil {
		l.Error("failed to insert scan summary", zap.Error(err))
		t.Fatalf("failed to insert scan summary: %v", err)
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

func TestUpsertScanSummary_Error(t *testing.T) {
	l := logger.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Expect transaction start
	mock.ExpectBegin()

	summary := models.ScanSummaryDAO{
		TotalVulnerabilities: 5,
		SeverityCounts: models.SeverityCountsDAO{
			Critical: 2,
			High:     1,
			Medium:   1,
			Low:      1,
		},
		FixableCount: 4,
		Compliant:    true,
	}

	// Expect INSERT query but return an error
	mock.ExpectExec("INSERT INTO scan_summary").
		WithArgs("scan_test", summary.TotalVulnerabilities, sqlmock.AnyArg(), summary.FixableCount, summary.Compliant).
		WillReturnError(fmt.Errorf("mock insert failure"))

	// Expect transaction rollback
	mock.ExpectRollback()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to start transaction: %v", err)
	}

	l.Info("starting error case for UpsertScanSummary test")

	err = UpsertScanSummary(tx, "scan_test", summary)
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