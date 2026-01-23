-- +goose Up
-- +goose StatementBegin

-- Add recipient_type column
ALTER TABLE bagsy_notifications ADD COLUMN recipient_type TEXT NOT NULL DEFAULT 'client';

-- Drop old unique constraint
ALTER TABLE bagsy_notifications DROP CONSTRAINT IF EXISTS bagsy_notifications_type_unique;

-- Add new unique constraint including recipient_type
ALTER TABLE bagsy_notifications ADD CONSTRAINT bagsy_notifications_unique UNIQUE(bagsy_id, type, recipient_type);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE bagsy_notifications DROP CONSTRAINT IF EXISTS bagsy_notifications_unique;
ALTER TABLE bagsy_notifications ADD CONSTRAINT bagsy_notifications_type_unique UNIQUE(bagsy_id, type);
ALTER TABLE bagsy_notifications DROP COLUMN recipient_type;

-- +goose StatementEnd
