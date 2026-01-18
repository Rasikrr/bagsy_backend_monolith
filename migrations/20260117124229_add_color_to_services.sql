-- +goose Up
-- +goose StatementBegin

ALTER TABLE services ADD COLUMN color SMALLINT NOT NULL DEFAULT 0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE services DROP COLUMN color;

-- +goose StatementEnd
