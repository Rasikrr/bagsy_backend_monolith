-- +goose Up
-- +goose StatementBegin
ALTER TABLE bagsies
    ADD COLUMN deleted_at TIMESTAMPTZ;

-- Создаем индекс для быстрых запросов с фильтром deleted_at IS NULL
CREATE INDEX idx_bagsies_deleted_at ON bagsies(deleted_at) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_bagsies_deleted_at;

ALTER TABLE bagsies
    DROP COLUMN deleted_at;
-- +goose StatementEnd
