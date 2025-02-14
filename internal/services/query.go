package services

import (
	"fmt"
	"kai-sec/internal/daos"
	"kai-sec/internal/dtos"
	"kai-sec/internal/logger"

	"go.uber.org/zap"
)


func GetVulnerabilities(filters dtos.FliterRequest) ([]dtos.Vulnerability, error) {
	l := logger.GetLogger()
	if filters.Filters.Severity == "" {
		l.Error("Missing severity filter")
		return nil, fmt.Errorf("severity filter is required")
	}

	l.Info("Fetching vulnerabilities", zap.String("severity", filters.Filters.Severity))

	// Fetch vulnerabilities from DB
	vulns, err := daos.GetVulnerabilitiesBySeverity(filters.Filters.Severity)
	if err != nil {
		l.Error("Error retrieving vulnerabilities", zap.Error(err))
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
