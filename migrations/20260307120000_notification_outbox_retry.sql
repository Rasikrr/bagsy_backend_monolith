-- +goose Up
-- +goose StatementBegin

ALTER TABLE notification_outbox
    ADD COLUMN status TEXT NOT NULL DEFAULT 'pending',
    ADD COLUMN recipient_phone TEXT NOT NULL DEFAULT '',
    ADD COLUMN recipient_type TEXT NOT NULL DEFAULT '',
    ADD COLUMN message TEXT NOT NULL DEFAULT '',
    ADD COLUMN appointment_id UUID,
    ADD COLUMN attempts INT NOT NULL DEFAULT 0,
    ADD COLUMN max_attempts INT NOT NULL DEFAULT 3,
    ADD COLUMN last_error TEXT,
    ADD COLUMN updated_at TIMESTAMPTZ;

ALTER TABLE notification_outbox
    DROP COLUMN IF EXISTS payload,
    DROP COLUMN IF EXISTS entity_id;

CREATE INDEX idx_notification_outbox_poll
    ON notification_outbox (status, scheduled_for)
    WHERE status = 'pending';

CREATE INDEX idx_notification_outbox_appointment
    ON notification_outbox (appointment_id)
    WHERE status = 'pending';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_notification_outbox_appointment;
DROP INDEX IF EXISTS idx_notification_outbox_poll;

ALTER TABLE notification_outbox
    ADD COLUMN IF NOT EXISTS payload JSONB NOT NULL DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS entity_id TEXT NOT NULL DEFAULT '';

ALTER TABLE notification_outbox
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS last_error,
    DROP COLUMN IF EXISTS max_attempts,
    DROP COLUMN IF EXISTS attempts,
    DROP COLUMN IF EXISTS appointment_id,
    DROP COLUMN IF EXISTS message,
    DROP COLUMN IF EXISTS recipient_type,
    DROP COLUMN IF EXISTS recipient_phone,
    DROP COLUMN IF EXISTS status;

-- +goose StatementEnd
