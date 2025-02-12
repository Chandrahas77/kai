package daos

import (
	"encoding/json"
	"fmt"
	"kai-sec/internal/daos/models"
	"kai-sec/internal/storage"
	"log"
)

func GetVulnerabilitiesBySeverity(severity string) ([]models.VulnerabilityDAO, error) {
	if storage.DB == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	query := `
		SELECT id, severity, cvss, status, package_name, current_version, fixed_version, 
		       description, published_date, link, risk_factors
		FROM vulnerabilities
		WHERE severity = $1
	`

	rows, err := storage.DB.Query(query, severity)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vulnerabilities: %w", err)
	}
	defer rows.Close()

	var vulnerabilities []models.VulnerabilityDAO

	for rows.Next() {
		var vuln models.VulnerabilityDAO
		var riskFactorsJSON string

		err := rows.Scan(
			&vuln.ID, &vuln.Severity, &vuln.CVSS, &vuln.Status, &vuln.PackageName,
			&vuln.CurrentVersion, &vuln.FixedVersion, &vuln.Description,
			&vuln.PublishedDate, &vuln.Link, &riskFactorsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert risk factors JSON string to slice
		vuln.RiskFactors = parseRiskFactors(riskFactorsJSON)
		vulnerabilities = append(vulnerabilities, vuln)
	}

	log.Printf("Fetched %d vulnerabilities with severity %s\n", len(vulnerabilities), severity)
	return vulnerabilities, nil
}

func parseRiskFactors(jsonString string) []string {
	var riskFactors []string
	if err := json.Unmarshal([]byte(jsonString), &riskFactors); err != nil {
		log.Printf("Error parsing risk factors JSON: %v\n", err)
		return []string{}
	}
	return riskFactors
}
