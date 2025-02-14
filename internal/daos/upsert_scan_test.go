package daos

import (
	"fmt"
	"kai-sec/internal/daos/models"
	"kai-sec/internal/storage"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUpsertScan_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()
	// Inject mock DB
	storage.DB = db 

	mock.ExpectBegin()

	//mock fetching
	mock.ExpectQuery("SELECT timestamp, scan_status, resource_type, resource_name FROM scans WHERE scan_id = \\$1").
		WithArgs("scan_001").
		WillReturnRows(sqlmock.NewRows([]string{"timestamp", "scan_status", "resource_type", "resource_name"}).
			AddRow(time.Now(), "completed", "container", "test-resource"))

	// mock scan insert/update
	mock.ExpectExec("INSERT INTO scans").
		WithArgs("scan_001", sqlmock.AnyArg(), "completed", "container", "test-resource").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// mock scan summary insert 
	mock.ExpectExec("INSERT INTO scan_summary").
		WithArgs("scan_001", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	//  mock scan metadata insert 
	mock.ExpectExec("INSERT INTO scan_metadata").
		WithArgs("scan_001", "", "", "[]", "[]"). //expecting empty lists as if it's nill
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	scan := models.ScanResultsDAO{
		ScanId:       "scan_001",
		Timestamp:    time.Now(),
		ScanStatus:   "completed",
		ResourceType: "container",
		ResourceName: "test-resource",
		ScanMetadata: models.ScanMetadataDAO{}, // Empty metadata
		Summary: models.ScanSummaryDAO{},       // Empty summary
	}

	err = UpsertScan(scan)
	if err != nil {
		t.Fatalf("UpsertScan failed: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet SQL expectations: %v", err)
	}
}

func TestUpsertScan_Failure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()
	storage.DB = db

	mock.ExpectBegin()

	
	mock.ExpectQuery("SELECT timestamp, scan_status, resource_type, resource_name FROM scans WHERE scan_id = \\$1").
		WithArgs("scan_001").
		WillReturnError(fmt.Errorf("database error"))

	mock.ExpectRollback()

	scan := models.ScanResultsDAO{
		ScanId:       "scan_001",
		Timestamp:    time.Now(),
		ScanStatus:   "completed",
		ResourceType: "container",
		ResourceName: "test-resource",
	}

	err = UpsertScan(scan)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet SQL expectations: %v", err)
	}
}