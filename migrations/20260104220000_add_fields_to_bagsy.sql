-- +goose Up
-- +goose StatementBegin
ALTER TABLE bagsies
    ADD COLUMN reject_reason TEXT,
    ADD COLUMN comment TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE bagsies
    DROP COLUMN reject_reason,
    DROP COLUMN comment;
-- +goose StatementEnd
