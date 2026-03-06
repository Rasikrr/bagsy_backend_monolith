-- +goose Up
-- +goose StatementBegin
ALTER TABLE notification_outbox
    ADD COLUMN metadata JSONB NOT NULL DEFAULT '{}',
    DROP COLUMN message;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE notification_outbox
    ADD COLUMN message TEXT NOT NULL DEFAULT '',
    DROP COLUMN metadata;
-- +goose StatementEnd
