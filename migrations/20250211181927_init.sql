-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS scans (
    scan_id TEXT PRIMARY KEY,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    scan_status TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS vulnerabilities (
    id TEXT PRIMARY KEY,
    scan_id TEXT NOT NULL,
    severity TEXT NOT NULL,
    cvss REAL NOT NULL,
    status TEXT NOT NULL,
    package_name TEXT NOT NULL,
    current_version TEXT NOT NULL,
    fixed_version TEXT NOT NULL,
    description TEXT NOT NULL,
    published_date TIMESTAMP NOT NULL,
    link TEXT NOT NULL,
    risk_factors JSONB NOT NULL,
    FOREIGN KEY (scan_id) REFERENCES scans(scan_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS scan_summary (
    scan_id TEXT PRIMARY KEY,
    total_vulnerabilities INT NOT NULL,
    severity_counts JSONB NOT NULL,
    fixable_count INT NOT NULL,
    compliant BOOLEAN NOT NULL,
    FOREIGN KEY (scan_id) REFERENCES scans(scan_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS scan_metadata (
    scan_id TEXT PRIMARY KEY,
    scanner_version TEXT NOT NULL,
    policies_version TEXT NOT NULL,
    scanning_rules JSONB NOT NULL,
    excluded_paths JSONB NOT NULL,
    FOREIGN KEY (scan_id) REFERENCES scans(scan_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS scan_metadata;
DROP TABLE IF EXISTS scan_summary;
DROP TABLE IF EXISTS vulnerabilities;
DROP TABLE IF EXISTS scans;
-- +goose StatementEnd
