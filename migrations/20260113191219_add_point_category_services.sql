-- +goose Up
-- +goose StatementBegin

-- Связь между категориями точек и категориями услуг (many-to-many)
CREATE TABLE IF NOT EXISTS point_category_services (
    id                   SERIAL PRIMARY KEY,
    point_category_id    INTEGER NOT NULL REFERENCES point_categories(id) ON DELETE CASCADE,
    service_category_id  INTEGER NOT NULL REFERENCES service_categories(id) ON DELETE CASCADE,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(point_category_id, service_category_id)
);

CREATE INDEX idx_pcs_point_category ON point_category_services(point_category_id);
CREATE INDEX idx_pcs_service_category ON point_category_services(service_category_id);

-- Seed начальные связи
INSERT INTO point_category_services (point_category_id, service_category_id) VALUES
-- Барбершоп (id=1): Стрижки, Бритье, Уход за бородой, Укладка, Окрашивание
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5),
-- Салон красоты (id=2): Стрижки, Укладка, Окрашивание, Маникюр, Педикюр, Косметология, Массаж
(2, 1), (2, 4), (2, 5), (2, 6), (2, 7), (2, 8), (2, 9),
-- Универсальный салон (id=3): все категории
(3, 1), (3, 2), (3, 3), (3, 4), (3, 5), (3, 6), (3, 7), (3, 8), (3, 9);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS point_category_services;

-- +goose StatementEnd
