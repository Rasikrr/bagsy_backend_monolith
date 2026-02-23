-- +goose Up
-- +goose StatementBegin

-- ═══════════════════════════════════════════════════════════════
-- Тестовые данные для флоу записи на услугу
-- ═══════════════════════════════════════════════════════════════
-- Владелец:      911f8256-c798-469f-9a78-07f43ada5679 (phone: 77715275251)
-- Организация:   edf875dd-dd42-408e-93b6-5606c13b74e0
-- Локация:       3322170f-ba5a-4799-baf8-ca86a8a0b22f

-- Обновляем локацию: schedule_type = mixed, шаг слота = 30 мин
UPDATE locations
SET schedule_type = 'mixed',
    slot_duration_minutes = 30
WHERE id = '3322170f-ba5a-4799-baf8-ca86a8a0b22f';

-- Разрешить владельцу оказывать услуги
UPDATE employees
SET can_provide_services = true
WHERE id = '911f8256-c798-469f-9a78-07f43ada5679';

-- ─────────────────────────────────────────────────────────────────
-- 1. Категории услуг (привязаны к location_category локации)
-- ─────────────────────────────────────────────────────────────────

-- Узнаём category_id локации и создаём service_categories
INSERT INTO service_categories (id, location_category_id, name, sort_order)
SELECT
    '11111111-0001-4000-8000-000000000001'::uuid,
    l.category_id,
    'Стрижки',
    10
FROM locations l WHERE l.id = '3322170f-ba5a-4799-baf8-ca86a8a0b22f';

INSERT INTO service_categories (id, location_category_id, name, sort_order)
SELECT
    '11111111-0001-4000-8000-000000000002'::uuid,
    l.category_id,
    'Окрашивание',
    20
FROM locations l WHERE l.id = '3322170f-ba5a-4799-baf8-ca86a8a0b22f';

INSERT INTO service_categories (id, location_category_id, name, sort_order)
SELECT
    '11111111-0001-4000-8000-000000000003'::uuid,
    l.category_id,
    'Уход за волосами',
    30
FROM locations l WHERE l.id = '3322170f-ba5a-4799-baf8-ca86a8a0b22f';

-- ─────────────────────────────────────────────────────────────────
-- 2. Услуги
-- ─────────────────────────────────────────────────────────────────

INSERT INTO services (id, location_id, category_id, name, description, duration_minutes, color, sort_order) VALUES
    ('22222222-0001-4000-8000-000000000001', '3322170f-ba5a-4799-baf8-ca86a8a0b22f',
     '11111111-0001-4000-8000-000000000001', 'Мужская стрижка', 'Классическая мужская стрижка', 30, 'green', 10),

    ('22222222-0001-4000-8000-000000000002', '3322170f-ba5a-4799-baf8-ca86a8a0b22f',
     '11111111-0001-4000-8000-000000000001', 'Женская стрижка', 'Стрижка любой длины', 60, 'purple', 20),

    ('22222222-0001-4000-8000-000000000003', '3322170f-ba5a-4799-baf8-ca86a8a0b22f',
     '11111111-0001-4000-8000-000000000002', 'Окрашивание корней', 'Окрашивание отросших корней', 90, 'red', 10),

    ('22222222-0001-4000-8000-000000000004', '3322170f-ba5a-4799-baf8-ca86a8a0b22f',
     '11111111-0001-4000-8000-000000000003', 'Кератиновое выпрямление', 'Восстановление и выпрямление волос', 120, 'orange', 10);

-- ─────────────────────────────────────────────────────────────────
-- 3. Привязка владельца к услугам (employee_services)
-- ─────────────────────────────────────────────────────────────────

INSERT INTO employee_services (employee_id, service_id, price) VALUES
    ('911f8256-c798-469f-9a78-07f43ada5679', '22222222-0001-4000-8000-000000000001', 3500.0000),  -- Мужская стрижка — 3500₸
    ('911f8256-c798-469f-9a78-07f43ada5679', '22222222-0001-4000-8000-000000000002', 5000.0000),  -- Женская стрижка — 5000₸
    ('911f8256-c798-469f-9a78-07f43ada5679', '22222222-0001-4000-8000-000000000003', 8000.0000),  -- Окрашивание — 8000₸
    ('911f8256-c798-469f-9a78-07f43ada5679', '22222222-0001-4000-8000-000000000004', 15000.0000); -- Кератин — 15000₸

-- ─────────────────────────────────────────────────────────────────
-- 4. Расписание локации на 30 дней (09:00–13:00 work, 13:00–14:00 rest, 14:00–21:00 work)
--    Воскресенье (dow=0) — выходной
-- ─────────────────────────────────────────────────────────────────

INSERT INTO location_schedules (location_id, date, type, start_time, end_time)
SELECT
    '3322170f-ba5a-4799-baf8-ca86a8a0b22f',
    d::date,
    slot.type,
    slot.start_time,
    slot.end_time
FROM generate_series(CURRENT_DATE, CURRENT_DATE + INTERVAL '30 days', '1 day') AS d
CROSS JOIN (
    VALUES
        ('work'::varchar, '09:00'::time, '13:00'::time),
        ('rest'::varchar, '13:00'::time, '14:00'::time),
        ('work'::varchar, '14:00'::time, '21:00'::time)
) AS slot(type, start_time, end_time)
WHERE EXTRACT(DOW FROM d) != 0;  -- без воскресенья

-- ─────────────────────────────────────────────────────────────────
-- 5. Расписание сотрудника (владельца) на 30 дней (10:00–15:00, rest 15:00–15:30, 15:30–19:00)
--    Воскресенье и понедельник — выходные
-- ─────────────────────────────────────────────────────────────────

INSERT INTO employee_schedules (employee_id, date, type, start_time, end_time)
SELECT
    '911f8256-c798-469f-9a78-07f43ada5679',
    d::date,
    slot.type,
    slot.start_time,
    slot.end_time
FROM generate_series(CURRENT_DATE, CURRENT_DATE + INTERVAL '30 days', '1 day') AS d
CROSS JOIN (
    VALUES
        ('work'::varchar, '10:00'::time, '15:00'::time),
        ('rest'::varchar, '15:00'::time, '15:30'::time),
        ('work'::varchar, '15:30'::time, '19:00'::time)
) AS slot(type, start_time, end_time)
WHERE EXTRACT(DOW FROM d) NOT IN (0, 1);  -- без воскресенья и понедельника

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM employee_schedules WHERE employee_id = '911f8256-c798-469f-9a78-07f43ada5679';

DELETE FROM location_schedules WHERE location_id = '3322170f-ba5a-4799-baf8-ca86a8a0b22f';

DELETE FROM employee_services WHERE employee_id = '911f8256-c798-469f-9a78-07f43ada5679'
    AND service_id IN (
        '22222222-0001-4000-8000-000000000001',
        '22222222-0001-4000-8000-000000000002',
        '22222222-0001-4000-8000-000000000003',
        '22222222-0001-4000-8000-000000000004'
    );

DELETE FROM services WHERE id IN (
    '22222222-0001-4000-8000-000000000001',
    '22222222-0001-4000-8000-000000000002',
    '22222222-0001-4000-8000-000000000003',
    '22222222-0001-4000-8000-000000000004'
);

DELETE FROM service_categories WHERE id IN (
    '11111111-0001-4000-8000-000000000001',
    '11111111-0001-4000-8000-000000000002',
    '11111111-0001-4000-8000-000000000003'
);

UPDATE employees SET can_provide_services = false WHERE id = '911f8256-c798-469f-9a78-07f43ada5679';
UPDATE locations SET schedule_type = NULL, slot_duration_minutes = 30 WHERE id = '3322170f-ba5a-4799-baf8-ca86a8a0b22f';

-- +goose StatementEnd
