package daos

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"kai-sec/internal/daos/models"
	"kai-sec/internal/storage"
	"log"
	"time"
)

func UpsertScan(scan models.ScanResultsDAO) error {
	if storage.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Starting a transaction 
	tx, err := storage.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = UpsertScanTable(tx, scan)
	if err != nil {
		return err
	}

	err = UpsertVulnerabilities(tx, scan.ScanId, scan.Vulnerabilities)
	if err != nil {
		return err
	}

	err = UpsertScanSummary(tx, scan.ScanId, scan.Summary)
	if err != nil {
		return err
	}

	err = UpsertScanMetadata(tx, scan.ScanId, scan.ScanMetadata)
	if err != nil {
		return err
	}

	log.Println("Successfully inserted/updated scan data.")
	return nil
}

func UpsertScanTable(tx *sql.Tx, scan models.ScanResultsDAO) error {
	// Check if the scan exists and whether any changes are needed
	var existingScan models.ScanResultsDAO
	err := tx.QueryRow(`
		SELECT timestamp, scan_status, resource_type, resource_name
		FROM scans WHERE scan_id = $1`, scan.ScanId).Scan(
		&existingScan.Timestamp, &existingScan.ScanStatus,
		&existingScan.ResourceType, &existingScan.ResourceName,
	)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing scan: %w", err)
	}

	// If the scan exists and the data is unchanged, only update `last_attempted_at`
	if err == nil && existingScan.Timestamp == scan.Timestamp &&
		existingScan.ScanStatus == scan.ScanStatus &&
		existingScan.ResourceType == scan.ResourceType &&
		existingScan.ResourceName == scan.ResourceName {
		var now time.Time
		_, err = tx.Exec(`UPDATE scans SET last_attempted_at = $1 WHERE scan_id = $2`, now, scan.ScanId)
		if err != nil {
			return fmt.Errorf("failed to update last_attempted_at: %w", err)
		}
		log.Printf("Scan ID %s already exists with the same data. Only updated last_attempted_at.\n", scan.ScanId)
		// Exit soon as there are no further updates needed
		return nil
	}

	//updating the scan as usual
	_, err = tx.Exec(`
		INSERT INTO scans (scan_id, timestamp, scan_status, resource_type, resource_name, created_at, updated_at, last_attempted_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW(), NOW())
		ON CONFLICT (scan_id) 
		DO UPDATE SET 
			timestamp = EXCLUDED.timestamp,
			scan_status = EXCLUDED.scan_status,
			resource_type = EXCLUDED.resource_type,
			resource_name = EXCLUDED.resource_name,
			updated_at = NOW(),
			last_attempted_at = NOW()`,
		scan.ScanId, scan.Timestamp, scan.ScanStatus, scan.ResourceType, scan.ResourceName,
	)
	if err != nil {
		return fmt.Errorf("failed to insert/update scan: %w", err)
	}
	log.Printf("Inserted or updated scan_id: %s\n", scan.ScanId)
	return nil
}

func UpsertVulnerabilities(tx *sql.Tx, scanId string, vulnerabilities []models.VulnerabilityDAO) error {
	for _, vuln := range vulnerabilities {
		riskFactorsJSON, _ := json.Marshal(vuln.RiskFactors)
		_, err := tx.Exec(`
			INSERT INTO vulnerabilities (id, scan_id, severity, cvss, status, package_name, 
				current_version, fixed_version, description, published_date, link, risk_factors, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
			ON CONFLICT (id) 
			DO UPDATE SET 
				severity = EXCLUDED.severity,
				cvss = EXCLUDED.cvss,
				status = EXCLUDED.status,
				package_name = EXCLUDED.package_name,
				current_version = EXCLUDED.current_version,
				fixed_version = EXCLUDED.fixed_version,
				description = EXCLUDED.description,
				published_date = EXCLUDED.published_date,
				link = EXCLUDED.link,
				risk_factors = EXCLUDED.risk_factors,
				updated_at = NOW()`,
			vuln.ID, scanId, vuln.Severity, vuln.CVSS, vuln.Status, vuln.PackageName,
			vuln.CurrentVersion, vuln.FixedVersion, vuln.Description, vuln.PublishedDate, vuln.Link, string(riskFactorsJSON),
		)
		if err != nil {
			log.Printf("Skipping vulnerability %s due to error: %v", vuln.ID, err)
		}
	}
	return nil
}

func UpsertScanMetadata(tx *sql.Tx, scanId string, metadata models.ScanMetadataDAO) error {
	scanningRulesJSON, _ := json.Marshal(metadata.ScanningRules)
	excludedPathsJSON, _ := json.Marshal(metadata.ExcludedPaths)

	_, err := tx.Exec(`
		INSERT INTO scan_metadata (scan_id, scanner_version, policies_version, scanning_rules, excluded_paths, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (scan_id) 
		DO UPDATE SET 
			scanner_version = EXCLUDED.scanner_version,
			policies_version = EXCLUDED.policies_version,
			scanning_rules = EXCLUDED.scanning_rules,
			excluded_paths = EXCLUDED.excluded_paths,
			updated_at = NOW()`,
		scanId, metadata.ScannerVersion, metadata.PoliciesVersion, string(scanningRulesJSON), string(excludedPathsJSON),
	)
	if err != nil {
		log.Printf("Skipping scan metadata for scan_id %s due to error: %v", scanId, err)
	}
	return nil
}

func UpsertScanSummary(tx *sql.Tx, scanId string, summary models.ScanSummaryDAO) error {
	severityCountsJSON, _ := json.Marshal(summary.SeverityCounts)

	_, err := tx.Exec(`
		INSERT INTO scan_summary (scan_id, total_vulnerabilities, severity_counts, fixable_count, compliant, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (scan_id) 
		DO UPDATE SET 
			total_vulnerabilities = EXCLUDED.total_vulnerabilities,
			severity_counts = EXCLUDED.severity_counts,
			fixable_count = EXCLUDED.fixable_count,
			compliant = EXCLUDED.compliant,
			updated_at = NOW()`,
		scanId, summary.TotalVulnerabilities, string(severityCountsJSON), summary.FixableCount, summary.Compliant,
	)
	if err != nil {
		log.Printf("Skipping scan summary for scan_id %s due to error: %v", scanId, err)
	}
	return nil
}

