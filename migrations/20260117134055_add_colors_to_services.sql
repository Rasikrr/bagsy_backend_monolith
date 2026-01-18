-- +goose Up
-- +goose StatementBegin

-- Назначаем цвета услугам по категориям:
-- 0 = black, 1 = green, 2 = red, 3 = yellow, 4 = purple, 5 = orange, 6 = gray

-- Стрижки (category_id=1) → Green (1)
UPDATE services SET color = 1 WHERE category_id = 1;

-- Бритье (category_id=2) → Red (2)
UPDATE services SET color = 2 WHERE category_id = 2;

-- Уход за бородой (category_id=3) → Orange (5)
UPDATE services SET color = 5 WHERE category_id = 3;

-- Укладка и стайлинг (category_id=4) → Yellow (3)
UPDATE services SET color = 3 WHERE category_id = 4;

-- Окрашивание волос (category_id=5) → Purple (4)
UPDATE services SET color = 4 WHERE category_id = 5;

-- Маникюр (category_id=6) → Green (1)
UPDATE services SET color = 1 WHERE category_id = 6;

-- Педикюр (category_id=7) → Red (2)
UPDATE services SET color = 2 WHERE category_id = 7;

-- Косметология (category_id=8) → Purple (4)
UPDATE services SET color = 4 WHERE category_id = 8;

-- Массаж (category_id=9) → Yellow (3)
UPDATE services SET color = 3 WHERE category_id = 9;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Сбрасываем все цвета на дефолтный (black = 0)
UPDATE services SET color = 0;

-- +goose StatementEnd
