package daos

import (
	"encoding/json"
	"fmt"
	"kai-sec/internal/daos/models"
	"kai-sec/internal/storage"
	"log"
)

func InsertScan(scan models.ScanResultsDAO) error {
	tx, err := storage.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	_, err = tx.Exec(`
		INSERT INTO scans (scan_id, timestamp, scan_status, resource_type, resource_name)
		VALUES ($1, $2, $3, $4, $5)`,
		scan.ScanId, scan.Timestamp, scan.ScanStatus, scan.ResourceType, scan.ResourceName,
	)
	if err != nil {
		return fmt.Errorf("failed to insert scan: %w", err)
	}

	for _, vuln := range scan.Vulnerabilities {
		riskFactorsJSON, _ := json.Marshal(vuln.RiskFactors)
		_, err = tx.Exec(`
			INSERT INTO vulnerabilities (id, scan_id, severity, cvss, status, package_name, 
				current_version, fixed_version, description, published_date, link, risk_factors)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			vuln.ID, scan.ScanId, vuln.Severity, vuln.CVSS, vuln.Status, vuln.PackageName,
			vuln.CurrentVersion, vuln.FixedVersion, vuln.Description, vuln.PublishedDate, vuln.Link, string(riskFactorsJSON),
		)
		if err != nil {
			return fmt.Errorf("failed to insert vulnerability: %w", err)
		}
	}

	severityCountsJSON, _ := json.Marshal(scan.Summary.SeverityCounts)
	_, err = tx.Exec(`
		INSERT INTO scan_summary (scan_id, total_vulnerabilities, severity_counts, fixable_count, compliant)
		VALUES ($1, $2, $3, $4, $5)`,
		scan.ScanId, scan.Summary.TotalVulnerabilities, string(severityCountsJSON), scan.Summary.FixableCount, scan.Summary.Compliant,
	)
	if err != nil {
		return fmt.Errorf("failed to insert scan summary: %w", err)
	}


	scanningRulesJSON, _ := json.Marshal(scan.ScanMetadata.ScanningRules)
	excludedPathsJSON, _ := json.Marshal(scan.ScanMetadata.ExcludedPaths)

	_, err = tx.Exec(`
		INSERT INTO scan_metadata (scan_id, scanner_version, policies_version, scanning_rules, excluded_paths)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (scan_id) DO NOTHING`,
		scan.ScanId, scan.ScanMetadata.ScannerVersion, scan.ScanMetadata.PoliciesVersion, string(scanningRulesJSON), string(excludedPathsJSON),
	)
	if err != nil {
		return fmt.Errorf("failed to insert scan metadata: %w", err)
	}

	// Commit transaction if everything succeeds
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Println("Successfully inserted scan data into PostgreSQL!")
	return nil
}
