package daos

import (
	"kai-sec/internal/daos/models"
	"kai-sec/internal/logger"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
)

func TestUpsertScanTable(t *testing.T) {
	l := logger.GetLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		l.Fatal("Failed to create mock DB", zap.Error(err))
	}
	defer db.Close()

	// Expect transaction start
	mock.ExpectBegin()

	// Expect a SELECT query for checking existing scan
	mock.ExpectQuery("SELECT timestamp, scan_status, resource_type, resource_name FROM scans WHERE scan_id = \\$1").
		WithArgs("scan_001").
		WillReturnRows(sqlmock.NewRows([]string{"timestamp", "scan_status", "resource_type", "resource_name"}).
			AddRow(time.Now(), "completed", "container", "test-resource"))

	// Expect INSERT INTO scans query
	mock.ExpectExec("INSERT INTO scans").
		WithArgs("scan_001", sqlmock.AnyArg(), "completed", "container", "test-resource").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect transaction commit
	mock.ExpectCommit()

	// Start a transaction in mock DB
	tx, err := db.Begin()
	if err != nil {
		l.Fatal("Failed to start transaction", zap.Error(err))
	}

	// Test scan data
	scan := models.ScanResultsDAO{
		ScanId:       "scan_001",
		Timestamp:    time.Now(),
		ScanStatus:   "completed",
		ResourceType: "container",
		ResourceName: "test-resource",
	}

	// Call UpsertScanTable with mock DB
	err = UpsertScanTable(tx, scan)
	if err != nil {
		t.Errorf("UpsertScanTable failed: %v", err)
	}

	// Commit transaction explicitly
	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Transaction expectations not met: %v", err)
	}
}

func TestUpsertScanTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Expect transaction start
	mock.ExpectBegin()

	// Expect SELECT query before upsert
	mock.ExpectQuery("SELECT timestamp, scan_status, resource_type, resource_name FROM scans WHERE scan_id = \\$1").
		WithArgs("scan_123").
		WillReturnRows(sqlmock.NewRows([]string{"timestamp", "scan_status", "resource_type", "resource_name"}).
			AddRow(time.Now(), "completed", "container", "test-resource"))

	// Expect INSERT INTO scans
	mock.ExpectExec("INSERT INTO scans").
		WithArgs("scan_123", sqlmock.AnyArg(), "completed", "container", "test-resource").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect transaction commit
	mock.ExpectCommit()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to start transaction: %v", err)
	}

	scan := models.ScanResultsDAO{
		ScanId:       "scan_123",
		Timestamp:    time.Now(),
		ScanStatus:   "completed",
		ResourceType: "container",
		ResourceName: "test-resource",
	}

	// Call UpsertScanTable
	err = UpsertScanTable(tx, scan)
	if err != nil {
		t.Errorf("upsertscantable transaction failed: %v", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		t.Fatalf("failed to commit transaction: %v", err)
	}

	// Ensure all expectations met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("transaction expectations not met: %v", err)
	}
}
