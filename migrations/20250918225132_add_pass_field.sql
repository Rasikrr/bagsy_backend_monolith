-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN password TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN password;
-- +goose StatementEnd
