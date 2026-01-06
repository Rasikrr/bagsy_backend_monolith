-- +goose Up
-- +goose StatementBegin

-- 1. Обновляем первичные ключи (networks.code и points.code должны быть обновлены первыми)
UPDATE networks SET code = LOWER(code) WHERE code != LOWER(code);
UPDATE points SET code = LOWER(code) WHERE code != LOWER(code);

-- 2. Обновляем внешние ключи в points (network_code)
UPDATE points SET network_code = LOWER(network_code) WHERE network_code != LOWER(network_code);

-- 3. Обновляем внешние ключи в users
UPDATE users SET point_code = LOWER(point_code) WHERE point_code IS NOT NULL AND point_code != LOWER(point_code);
UPDATE users SET network_code = LOWER(network_code) WHERE network_code IS NOT NULL AND network_code != LOWER(network_code);

-- 4. Обновляем внешние ключи в services
UPDATE services SET point_code = LOWER(point_code) WHERE point_code != LOWER(point_code);

-- 5. Обновляем внешние ключи в bagsies
UPDATE bagsies SET point_code = LOWER(point_code) WHERE point_code != LOWER(point_code);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Откат невозможен, т.к. мы не знаем исходный регистр символов
-- Можно было бы сохранить старые значения в временную таблицу, но это излишне
SELECT 'Cannot rollback - original case is lost';

-- +goose StatementEnd
