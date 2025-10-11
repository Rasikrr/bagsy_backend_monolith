-- +goose Up
-- +goose StatementBegin

UPDATE users SET updated_by = 'system' WHERE updated_by IS NULL;

ALTER TABLE users
    ALTER COLUMN updated_by SET DEFAULT 'system',
    ALTER COLUMN updated_by SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE users
    ALTER COLUMN updated_by DROP DEFAULT,
    ALTER COLUMN updated_by DROP NOT NULL;
-- +goose StatementEnd
