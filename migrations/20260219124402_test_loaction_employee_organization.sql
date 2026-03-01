-- +goose Up
-- +goose StatementBegin

INSERT INTO organizations (
    id,
    name,
    description,
    slug,
    active,
    created_at,
    updated_at,
    deleted_at
) VALUES (
    'edf875dd-dd42-408e-93b6-5606c13b74e0', -- id
    '',                                     -- name (пусто на скрине)
    NULL,                                   -- description
    NULL,                                   -- slug
    true,                                   -- active
    '2026-02-23 16:19:23.757126+05',        -- created_at
    NULL,                                   -- updated_at
    NULL                                    -- deleted_at
);

INSERT INTO locations (
    id,
    organization_id,
    category_id,
    name,
    description,
    phone,
    slug,
    city,
    address_street,
    address_building,
    address_details,
    longitude,
    latitude,
    active,
    schedule_type,
    slot_duration_minutes,
    created_at,
    updated_at,
    deleted_at
) VALUES (
    '3322170f-ba5a-4799-baf8-ca86a8a0b22f', -- id
    'edf875dd-dd42-408e-93b6-5606c13b74e0', -- organization_id
    (SELECT id FROM location_categories lc WHERE lc.slug = 'beauty-salon'), -- category_id
    'NAVI',                                 -- name
    'Премиум салон красоты',                -- description
    '77715275251',                          -- phone
    'navi',                                 -- slug
    'Астана',                               -- city
    'Кунаева',                              -- address_street
    '68',                                   -- address_building
    NULL,                                   -- address_details (пусто на скрине)
    73.563,                                 -- longitude
    74.7545,                                -- latitude
    true,                                   -- active
    'mixed',                                -- schedule_type
    30,                                     -- slot_duration_minutes
    '2026-02-23 16:28:31.570866+05',        -- created_at
    NULL,                                   -- updated_at
    NULL                                    -- deleted_at
);



INSERT INTO employees (
    id,
    phone,
    password_hash,
    first_name,
    last_name,
    organization_id,
    location_id,
    role,
    can_provide_services,
    can_manage_location_schedule,
    active,
    created_at,
    updated_at,
    deleted_at,
    avatar_id
) VALUES (
    '911f8256-c798-469f-9a78-07f43ada5679', -- id
    '77715275251',                          -- phone
    '$2a$10$gmd4gREMdZrPW6L5NXcj2emx8jZEj1xL9JrrUfDGjf/S..qx6b/k.', -- password_hash
    'Расул',                                -- first_name
    'Туртулов',                             -- last_name
    'edf875dd-dd42-408e-93b6-5606c13b74e0', -- organization_id
    '3322170f-ba5a-4799-baf8-ca86a8a0b22f', -- location_id
    'owner',                                -- role
    true,                                   -- can_provide_services
    true,                                   -- can_manage_location_schedule
    true,                                   -- active
    '2026-02-23 16:19:23.759124+05',        -- created_at
    '2026-02-23 16:28:31.573319+05',        -- updated_at
    NULL,                                   -- deleted_at
    NULL                                    -- avatar_id
);


-- Подписка (plan: point, status: active)
INSERT INTO subscriptions (
    id, organization_id, plan_id, status, billing_cycle,
    recurring_amount, current_period_start, current_period_end
) VALUES (
     'b0b0b0b0-0000-4000-8000-000000000001',
     'edf875dd-dd42-408e-93b6-5606c13b74e0',
     (SELECT id FROM plans WHERE code = 'point'),
     'active',
     'monthly',
     25000.00,
     NOW(),
     NOW() + INTERVAL '30 days'
 );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
