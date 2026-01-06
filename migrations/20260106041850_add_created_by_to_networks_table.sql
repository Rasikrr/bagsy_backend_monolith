-- +goose Up
-- +goose StatementBegin
ALTER TABLE networks
    ADD COLUMN created_by TEXT NOT NULL DEFAULT 'system';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE networks
    DROP COLUMN created_by;
-- +goose StatementEnd
