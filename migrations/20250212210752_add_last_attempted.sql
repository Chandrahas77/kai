-- +goose Up
-- +goose StatementBegin
ALTER TABLE scans ADD COLUMN IF NOT EXISTS last_attempted_at TIMESTAMP DEFAULT NOW();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE scans DROP COLUMN IF EXISTS last_attempted_at;
-- +goose StatementEnd
