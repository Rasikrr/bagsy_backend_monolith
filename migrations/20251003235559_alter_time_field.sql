-- +goose Up
-- +goose StatementBegin
ALTER TABLE bagsies
DROP
COLUMN time;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE bagsies
    ADD COLUMN time TIMESTAMP NOT NULL;
-- +goose StatementEnd
