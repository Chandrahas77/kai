package daos

import (
	"kai-sec/internal/storage"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetVulnerabilitiesBySeverity(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()
	storage.DB = db

	mock.ExpectQuery("SELECT id, severity, cvss, status, package_name, current_version, fixed_version, description, published_date, link, risk_factors FROM vulnerabilities WHERE severity = \\$1").
		WithArgs("HIGH").
		WillReturnRows(sqlmock.NewRows([]string{"id", "severity", "cvss", "status", "package_name", "current_version", "fixed_version", "description", "published_date", "link", "risk_factors"}).
			AddRow("CVE-2024-9999", "HIGH", 8.5, "active", "openssl", "1.1.1t", "1.1.1u", "Buffer overflow", "2024-01-01", "https://example.com", "[]"))

	_, err = GetVulnerabilitiesBySeverity("HIGH")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}