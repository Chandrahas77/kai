package services

import (
	"fmt"
	"kai-sec/internal/daos"
	"kai-sec/internal/dtos"
)

func GetVulnerabilities(filters dtos.FliterRequest) ([]dtos.Vulnerability, error) {
	if filters.Filters.Severity == "" {
		return nil, fmt.Errorf("severity filter is required")
	}

	// Fetch vulnerabilities from DB
	vulns, err := daos.GetVulnerabilitiesBySeverity(filters.Filters.Severity)
	if err != nil {
		return nil, fmt.Errorf("error retrieving vulnerabilities: %w", err)
	}

	// Convert daso to dtos for response
	var results []dtos.Vulnerability
	for _, v := range vulns {
		results = append(results, dtos.Vulnerability{
			ID:             v.ID,
			Severity:       v.Severity,
			CVSS:           v.CVSS,
			Status:         v.Status,
			PackageName:    v.PackageName,
			CurrentVersion: v.CurrentVersion,
			FixedVersion:   v.FixedVersion,
			Description:    v.Description,
			PublishedDate:  v.PublishedDate,
			Link:           v.Link,
			RiskFactors:    v.RiskFactors,
		})
	}

	return results, nil
}
