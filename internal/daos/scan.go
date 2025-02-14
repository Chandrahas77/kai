package daos

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"kai-sec/internal/daos/models"
	"kai-sec/internal/logger"
	"kai-sec/internal/storage"
	"time"

	"go.uber.org/zap"
)

func UpsertScan(scan models.ScanResultsDAO) error {
	l := logger.GetLogger()
	if storage.DB == nil {
		err := fmt.Errorf("database connection is nil")
		l.Error("UpsertScan failed", zap.Error(err))
		return err
	}

	// starting a transaction
	tx, err := storage.DB.Begin()
	if err != nil {
		l.Error("Failed to begin transaction", zap.Error(err))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			l.Error("Transaction rolled back due to error", zap.Error(err))
		} else {
			err = tx.Commit()
			if err != nil {
				l.Error("Transaction commit failed", zap.Error(err))
			}
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

	l.Info("successfully inserted/updated scan data", zap.String("scan_id", scan.ScanId))
	return nil
}

func UpsertScanTable(exec ExecContext, scan models.ScanResultsDAO) error {
	l := logger.GetLogger()
	// look for existing scan to avoid duplicates
	var existingScan models.ScanResultsDAO
	err := exec.QueryRow(`
		SELECT timestamp, scan_status, resource_type, resource_name
		FROM scans WHERE scan_id = $1`, scan.ScanId).Scan(
		&existingScan.Timestamp, &existingScan.ScanStatus,
		&existingScan.ResourceType, &existingScan.ResourceName,
	)
	if err != nil && err != sql.ErrNoRows {
		l.Error("failed to check existing scan", zap.String("scan_id", scan.ScanId), zap.Error(err))
		return fmt.Errorf("failed to check existing scan: %w", err)
	}

	// If the scan exists and the data is unchanged, only update `last_attempted_at`
	if err == nil && existingScan.Timestamp == scan.Timestamp &&
		existingScan.ScanStatus == scan.ScanStatus &&
		existingScan.ResourceType == scan.ResourceType &&
		existingScan.ResourceName == scan.ResourceName {
		var now time.Time
		_, err = exec.Exec(`UPDATE scans SET last_attempted_at = $1 WHERE scan_id = $2`, now, scan.ScanId)
		if err != nil {
			l.Error("Failed to update last_attempted_at", zap.String("scan_id", scan.ScanId), zap.Error(err))
			return fmt.Errorf("failed to update last_attempted_at: %w", err)
		}
		l.Info("scan exists, updated last_attempted_at", zap.String("scan_id", scan.ScanId))
		return nil
	}

	//Upsert Scan
	_, err = exec.Exec(`
		INSERT INTO scans (scan_id, timestamp, scan_status, resource_type, resource_name, created_at, updated_at, last_attempted_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (scan_id)
		DO UPDATE SET
			timestamp = EXCLUDED.timestamp,
			scan_status = EXCLUDED.scan_status,
			resource_type = EXCLUDED.resource_type,
			resource_name = EXCLUDED.resource_name,
			updated_at = CURRENT_TIMESTAMP,
			last_attempted_at = CURRENT_TIMESTAMP`,
		scan.ScanId, scan.Timestamp, scan.ScanStatus, scan.ResourceType, scan.ResourceName,
	)
	if err != nil {
		l.Error("failed to upsert scan", zap.String("scan_id", scan.ScanId), zap.Error(err))
		return fmt.Errorf("failed to insert/update scan: %w", err)
	}
	l.Info("upserted scan", zap.String("scan_id", scan.ScanId))
	return nil
}

func UpsertVulnerabilities(exec ExecContext, scanId string, vulnerabilities []models.VulnerabilityDAO) error {
	l := logger.GetLogger()
	for _, vuln := range vulnerabilities {
		riskFactorsJSON, _ := json.Marshal(vuln.RiskFactors)
		_, err := exec.Exec(`
			INSERT INTO vulnerabilities (id, scan_id, severity, cvss, status, package_name, 
				current_version, fixed_version, description, published_date, link, risk_factors, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
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
				updated_at = CURRENT_TIMESTAMP`,
			vuln.ID, scanId, vuln.Severity, vuln.CVSS, vuln.Status, vuln.PackageName,
			vuln.CurrentVersion, vuln.FixedVersion, vuln.Description, vuln.PublishedDate, vuln.Link, string(riskFactorsJSON),
		)
		if err != nil {
			l.Error("failed to upsert vulnerability",
				zap.String("vulnerability_id", vuln.ID),
				zap.String("scan_id", scanId),
				zap.Error(err),
			)
			return err
		}
	}
	return nil
}

func UpsertScanMetadata(exec ExecContext, scanId string, metadata models.ScanMetadataDAO) error {
	l := logger.GetLogger()
	scanningRulesJSON, _ := json.Marshal(metadata.ScanningRules)
	if len(metadata.ScanningRules) == 0 {
		scanningRulesJSON = []byte("[]")
	}
	excludedPathsJSON, _ := json.Marshal(metadata.ExcludedPaths)
	if len(metadata.ExcludedPaths) == 0 {
		excludedPathsJSON = []byte("[]")
	}

	_, err := exec.Exec(`
		INSERT INTO scan_metadata (scan_id, scanner_version, policies_version, scanning_rules, excluded_paths, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (scan_id) 
		DO UPDATE SET 
			scanner_version = EXCLUDED.scanner_version,
			policies_version = EXCLUDED.policies_version,
			scanning_rules = EXCLUDED.scanning_rules,
			excluded_paths = EXCLUDED.excluded_paths,
			updated_at = CURRENT_TIMESTAMP`,
		scanId, metadata.ScannerVersion, metadata.PoliciesVersion, string(scanningRulesJSON), string(excludedPathsJSON),
	)
	if err != nil {
		l.Error("failed to upsert scan metadata", zap.String("scan_id", scanId), zap.Error(err))
		return err
	}
	l.Info("successfully upserted scan metadata", zap.String("scan_id", scanId))
	return nil
}

func UpsertScanSummary(exec ExecContext, scanId string, summary models.ScanSummaryDAO) error {
	l := logger.GetLogger()

	severityCountsJSON, _ := json.Marshal(summary.SeverityCounts)

	_, err := exec.Exec(`
		INSERT INTO scan_summary (scan_id, total_vulnerabilities, severity_counts, fixable_count, compliant, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (scan_id) 
		DO UPDATE SET 
			total_vulnerabilities = EXCLUDED.total_vulnerabilities,
			severity_counts = EXCLUDED.severity_counts,
			fixable_count = EXCLUDED.fixable_count,
			compliant = EXCLUDED.compliant,
			updated_at = CURRENT_TIMESTAMP`,
		scanId, summary.TotalVulnerabilities, string(severityCountsJSON), summary.FixableCount, summary.Compliant,
	)
	if err != nil {
		l.Error("skipping scan summary for scan_id %s due to error: %v", zap.String("scan_id", scanId), zap.Error(err))
		return err
	}
	return nil
}
