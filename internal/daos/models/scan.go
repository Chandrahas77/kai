package models

import "time"

type ScanReportDAO struct {
	ScanResults ScanResultsDAO `json:"scanResults"`
}

type ScanResultsDAO struct {
	ScanId          string             `json:"scan_id"`
	Timestamp        time.Time          `json:"timestamp"`
	ScanStatus      string             `json:"scan_status"`
	ResourceType    string             `json:"resource_type"`
	ResourceName    string             `json:"resource_name"`
	Vulnerabilities []VulnerabilityDAO `json:"vulnerabilities"`
	Summary         ScanSummaryDAO     `json:"summary"`
	ScanMetadata    ScanMetadataDAO    `json:"scan_metadata"`
}

type VulnerabilityDAO struct {
	ID             string    `json:"id"`
	Severity       string    `json:"severity"`
	CVSS           float64   `json:"cvss"`
	Status         string    `json:"status"`
	PackageName    string    `json:"package_name"`
	CurrentVersion string    `json:"current_version"`
	FixedVersion   string    `json:"fixed_version"`
	Description    string    `json:"description"`
	PublishedDate  time.Time `json:"published_date"`
	Link           string    `json:"link"`
	RiskFactors    []string  `json:"risk_factors"`
}

type ScanSummaryDAO struct {
	TotalVulnerabilities int               `json:"total_vulnerabilities"`
	SeverityCounts       SeverityCountsDAO `json:"severity_counts"`
	FixableCount         int               `json:"fixable_count"`
	Compliant            bool              `json:"complaint"`
}

type SeverityCountsDAO struct {
	Critical int `json:"CRITICAL"`
	High     int `json:"HIGH"`
	Medium   int `json:"MEDIUM"`
	Low      int `json:"LOW"`
}

type ScanMetadataDAO struct {
	ScannerVersion  string   `json:"scanner_version"`
	PoliciesVersion string   `json:"policies_version"`
	ScanningRules   []string `json:"scanning_rules"`
	ExcludedPaths   []string `json:"excluded_paths"`
}
