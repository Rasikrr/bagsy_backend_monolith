-- +goose Up
-- +goose StatementBegin

-- Включаем расширение pg_trgm для триграммного поиска
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- GIN-индекс для эффективного поиска по части номера телефона
-- Поддерживает поиск в начале, середине и конце строки
CREATE INDEX IF NOT EXISTS idx_users_phone_trgm ON users USING GIN (phone gin_trgm_ops);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_users_phone_trgm;
-- Не удаляем расширение pg_trgm, так как оно может использоваться другими таблицами

-- +goose StatementEnd
