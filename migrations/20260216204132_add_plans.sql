-- +goose Up
-- +goose StatementBegin
INSERT INTO plans (code, name, description, price_monthly, price_annual, sort_order, active)
VALUES
    ('solo',    'Solo',    'Для самозанятых мастеров',       5000.00, 48000.00,  0, true),
    ('point',   'Point',   'Для одной точки с сотрудниками', 9000.00, 86400.00,  1, true),
    ('network', 'Network', 'Для сети из нескольких точек',  25000.00, 240000.00, 2, true);

INSERT INTO plan_capabilities (plan_id, resource, limit_value) VALUES
    ((SELECT id FROM plans WHERE code = 'solo'), 'max_locations', 1),
    ((SELECT id FROM plans WHERE code = 'solo'), 'max_employees', 1),

    ((SELECT id FROM plans WHERE code = 'point'), 'max_locations', 1),
    ((SELECT id FROM plans WHERE code = 'point'), 'max_employees', 10),

    ((SELECT id FROM plans WHERE code = 'network'), 'max_locations', NULL),
    ((SELECT id FROM plans WHERE code = 'network'), 'max_employees', NULL)
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
