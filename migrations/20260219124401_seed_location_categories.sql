-- +goose Up
-- +goose StatementBegin
INSERT INTO location_categories (slug, name, sort_order) VALUES
    ('beauty-salon', 'Салон красоты', 10),
    ('barbershop', 'Барбершоп', 20),
    ('nail-studio', 'Студия маникюра', 30),
    ('spa-wellness', 'СПА и велнес', 40),
    ('dental-clinic', 'Стоматология', 50),
    ('medical-center', 'Медицинский центр', 60),
    ('fitness-club', 'Фитнес-клуб', 70),
    ('car-service', 'Автосервис', 80),
    ('pet-grooming', 'Груминг-салон', 90),
    ('educational-center', 'Учебный центр', 100);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE location_categories CASCADE;
-- +goose StatementEnd
