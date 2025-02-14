package daos

import (
	"encoding/json"
	"fmt"
	"kai-sec/internal/daos/models"
	"kai-sec/internal/logger"
	"kai-sec/internal/storage"
	"time"

	"go.uber.org/zap"
)

func GetVulnerabilitiesBySeverity(severity string) ([]models.VulnerabilityDAO, error) {
	l := logger.GetLogger()

	query := `
		SELECT id, severity, cvss, status, package_name, current_version, fixed_version, 
		       description, published_date, link, risk_factors
		FROM vulnerabilities
		WHERE severity = $1
	`

	rows, err := storage.DB.Query(query, severity)
	if err != nil {
		l.Error("Failed to fetch vulnerabilities", zap.String("severity", severity), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch vulnerabilities: %w", err)
	}
	defer rows.Close()

	var vulnerabilities []models.VulnerabilityDAO

	for rows.Next() {
		var vuln models.VulnerabilityDAO
		var riskFactorsJSON string
		var publishedDateStr string

		err := rows.Scan(
			&vuln.ID, &vuln.Severity, &vuln.CVSS, &vuln.Status, &vuln.PackageName,
			&vuln.CurrentVersion, &vuln.FixedVersion, &vuln.Description,
			&publishedDateStr, &vuln.Link, &riskFactorsJSON,
		)
		if err != nil {
			l.Error("Failed to scan row", zap.String("severity", severity), zap.Error(err))
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		vuln.PublishedDate, err = time.Parse(time.RFC3339, publishedDateStr)
		if err != nil {
			vuln.PublishedDate = time.Time{} 
		}

		// Convert risk factors JSON string to slice
		vuln.RiskFactors = parseRiskFactors(riskFactorsJSON)
		vulnerabilities = append(vulnerabilities, vuln)
	}

	l.Info("Fetched vulnerabilities successfully", zap.String("severity", severity), zap.Int("count", len(vulnerabilities)))
	return vulnerabilities, nil
}

func parseRiskFactors(jsonString string) []string {
	l := logger.GetLogger()
	var riskFactors []string
	if err := json.Unmarshal([]byte(jsonString), &riskFactors); err != nil {
		l.Warn("Error parsing risk factors JSON", zap.String("json", jsonString), zap.Error(err))
		return []string{}
	}
	return riskFactors
}
