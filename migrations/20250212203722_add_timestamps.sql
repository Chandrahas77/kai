-- +goose Up
-- +goose StatementBegin
ALTER TABLE scans ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT NOW();
ALTER TABLE scans ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();

ALTER TABLE vulnerabilities ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT NOW();
ALTER TABLE vulnerabilities ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();

ALTER TABLE scan_summary ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT NOW();
ALTER TABLE scan_summary ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();

ALTER TABLE scan_metadata ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT NOW();
ALTER TABLE scan_metadata ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE scans DROP COLUMN IF EXISTS created_at;
ALTER TABLE scans DROP COLUMN IF EXISTS updated_at;

ALTER TABLE vulnerabilities DROP COLUMN IF EXISTS created_at;
ALTER TABLE vulnerabilities DROP COLUMN IF EXISTS updated_at;

ALTER TABLE scan_summary DROP COLUMN IF EXISTS created_at;
ALTER TABLE scan_summary DROP COLUMN IF EXISTS updated_at;

ALTER TABLE scan_metadata DROP COLUMN IF EXISTS created_at;
ALTER TABLE scan_metadata DROP COLUMN IF EXISTS updated_at;
-- +goose StatementEnd
