-- +goose Up
-- +goose StatementBegin

-- Тестовые данные для GET /api/v1/employees
-- Организация: edf875dd-dd42-408e-93b6-5606c13b74e0
-- Локация:     3322170f-ba5a-4799-baf8-ca86a8a0b22f
-- Пароль у всех: password123


-- Обновить пароль owner-а на password123
UPDATE employees
SET password_hash = '$2b$12$EORmReKrYZdT3h658S1gMug/Le1S0ODueNC27PjF2hUA/7IyTh47e'
WHERE id = '911f8256-c798-469f-9a78-07f43ada5679';

-- Manager
INSERT INTO employees (
    id, phone, password_hash, first_name, last_name,
    organization_id, location_id, role,
    can_provide_services, can_manage_location_schedule,
    active, created_at, updated_at, deleted_at, avatar_id
) VALUES (
    'a1b2c3d4-1111-4000-8000-000000000001',
    '77001112233',
    '$2b$12$EORmReKrYZdT3h658S1gMug/Le1S0ODueNC27PjF2hUA/7IyTh47e',
    'Алия', 'Кенесова',
    'edf875dd-dd42-408e-93b6-5606c13b74e0',
    '3322170f-ba5a-4799-baf8-ca86a8a0b22f',
    'manager',
    false, true,
    true, NOW(), NULL, NULL, NULL
);

-- Staff (активный, может оказывать услуги)
INSERT INTO employees (
    id, phone, password_hash, first_name, last_name,
    organization_id, location_id, role,
    can_provide_services, can_manage_location_schedule,
    active, created_at, updated_at, deleted_at, avatar_id
) VALUES (
    'a1b2c3d4-2222-4000-8000-000000000002',
    '77002223344',
    '$2b$12$EORmReKrYZdT3h658S1gMug/Le1S0ODueNC27PjF2hUA/7IyTh47e',
    'Дамир', 'Ахметов',
    'edf875dd-dd42-408e-93b6-5606c13b74e0',
    '3322170f-ba5a-4799-baf8-ca86a8a0b22f',
    'staff',
    true, false,
    true, NOW(), NULL, NULL, NULL
);

-- Staff (активный, может оказывать услуги)
INSERT INTO employees (
    id, phone, password_hash, first_name, last_name,
    organization_id, location_id, role,
    can_provide_services, can_manage_location_schedule,
    active, created_at, updated_at, deleted_at, avatar_id
) VALUES (
    'a1b2c3d4-3333-4000-8000-000000000003',
    '77003334455',
    '$2b$12$EORmReKrYZdT3h658S1gMug/Le1S0ODueNC27PjF2hUA/7IyTh47e',
    'Айгерим', 'Сулейменова',
    'edf875dd-dd42-408e-93b6-5606c13b74e0',
    '3322170f-ba5a-4799-baf8-ca86a8a0b22f',
    'staff',
    true, false,
    true, NOW(), NULL, NULL, NULL
);

-- Staff (неактивный — для теста фильтра active=false)
INSERT INTO employees (
    id, phone, password_hash, first_name, last_name,
    organization_id, location_id, role,
    can_provide_services, can_manage_location_schedule,
    active, created_at, updated_at, deleted_at, avatar_id
) VALUES (
    'a1b2c3d4-4444-4000-8000-000000000004',
    '77004445566',
    '$2b$12$EORmReKrYZdT3h658S1gMug/Le1S0ODueNC27PjF2hUA/7IyTh47e',
    'Ержан', 'Бекмуратов',
    'edf875dd-dd42-408e-93b6-5606c13b74e0',
    '3322170f-ba5a-4799-baf8-ca86a8a0b22f',
    'staff',
    true, false,
    false, NOW(), NULL, NULL, NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM employees WHERE id IN (
    'a1b2c3d4-1111-4000-8000-000000000001',
    'a1b2c3d4-2222-4000-8000-000000000002',
    'a1b2c3d4-3333-4000-8000-000000000003',
    'a1b2c3d4-4444-4000-8000-000000000004'
);

DELETE FROM subscriptions WHERE id = 'b0b0b0b0-0000-4000-8000-000000000001';

-- +goose StatementEnd
