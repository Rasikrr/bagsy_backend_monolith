-- +goose Up
-- +goose StatementBegin
ALTER TABLE notification_outbox
    DROP COLUMN locked_until;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE notification_outbox
    ADD COLUMN locked_until TIMESTAMPTZ;
-- +goose StatementEnd
