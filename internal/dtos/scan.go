package dtos

import "time"

type ScanReport struct {
	ScanResults ScanResults `json:"scanResults"`
}

type ScanResults struct {
	ScanId          string          `json:"scan_id"`
	Timetamp        time.Time       `json:"timestamp"`
	ScanStatus      string          `json:"scan_status"`
	ResourceType    string          `json:"resource_type"`
	ResourceName    string          `json:"resource_name"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	Summary         ScanSummary     `json:"summary"`
	ScanMetadata    ScanMetadata    `json:"scan_metadata"`
}

type Vulnerability struct {
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

type ScanSummary struct {
	TotalVulnerabilities int            `json:"total_vulnerabilities"`
	SeverityCounts       SeverityCounts `json:"severity_counts"`
	FixableCount         int            `json:"fixable_count"`
	Complaint            bool           `json:"complaint"`
}

type SeverityCounts struct {
	Critical int `json:"CRITICAL"`
	High     int `json:"HIGH"`
	Medium   int `json:"MEDIUM"`
	Low      int `json:"LOW"`
}

type ScanMetadata struct {
	ScannerVersion  string   `json:"scanner_version"`
	PoliciesVersion string   `json:"policies_version"`
	ScanningRules   []string `json:"scanning_rules"`
	ExcludedPaths   []string `json:"excluded_paths"`
}
