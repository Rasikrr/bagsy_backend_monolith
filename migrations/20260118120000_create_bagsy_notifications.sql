-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS bagsy_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bagsy_id UUID NOT NULL REFERENCES bagsies(id) ON DELETE CASCADE,
    type TEXT NOT NULL,                           -- day_before, hour_before
    scheduled_at TIMESTAMPTZ NOT NULL,            -- когда должно быть отправлено
    sent_at TIMESTAMPTZ,                          -- когда фактически отправлено
    status TEXT NOT NULL DEFAULT 'pending',       -- pending, sent, failed, skipped
    attempts INT NOT NULL DEFAULT 0,
    last_error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT bagsy_notifications_type_unique UNIQUE(bagsy_id, type)
);

-- Индекс для поиска pending уведомлений по времени
CREATE INDEX idx_bagsy_notifications_pending_scheduled
    ON bagsy_notifications(scheduled_at)
    WHERE status = 'pending';

-- Индекс для поиска уведомлений по bagsy_id
CREATE INDEX idx_bagsy_notifications_bagsy_id ON bagsy_notifications(bagsy_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS bagsy_notifications;

-- +goose StatementEnd
