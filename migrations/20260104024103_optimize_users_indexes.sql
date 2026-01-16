-- +goose Up
-- +goose StatementBegin

-- 2. Фильтрация по сети/точке/роли + статус удаления + сортировка.
-- Порядок: Equality (network) -> Equality (deleted_at) -> Sort (created_at)
CREATE INDEX IF NOT EXISTS idx_users_network_deleted_created
    ON users(network_code, deleted_at, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_users_point_deleted_created
    ON users(point_code, deleted_at, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_users_role_deleted_created
    ON users(role, deleted_at, created_at DESC);

-- 3. Для общего списка (всех удаленных или всех активных)
CREATE INDEX IF NOT EXISTS idx_users_deleted_created
    ON users(deleted_at, created_at DESC);

-- 4. Для поиска по именам (с учетом удаления)
CREATE INDEX IF NOT EXISTS idx_users_names_deleted
    ON users(surname, name, deleted_at);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin

-- Удаляем новые индексы
DROP INDEX IF EXISTS idx_users_names_deleted;
DROP INDEX IF EXISTS idx_users_deleted_created;
DROP INDEX IF EXISTS idx_users_role_deleted_created;
DROP INDEX IF EXISTS idx_users_point_deleted_created;
DROP INDEX IF EXISTS idx_users_network_deleted_created;

-- Восстанавливаем старые индексы
CREATE INDEX IF NOT EXISTS idx_users_point_code ON users(point_code);
CREATE INDEX IF NOT EXISTS idx_users_network_code ON users(network_code);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- +goose StatementEnd
