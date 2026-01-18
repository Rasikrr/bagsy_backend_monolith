-- +goose Up
-- +goose StatementBegin

-- ============================================
-- ЗАПИСИ (BAGSIES) ДЛЯ ВСЕХ МАСТЕРОВ
-- ============================================
-- Для каждого staff добавляем:
-- - 5 записей на текущую неделю (13-19 января 2026)
-- - 3 записи на следующую неделю (20-26 января 2026)
--
-- Статусы: pending, created, completed, canceled
-- ============================================

-- Тестовые клиенты
-- 77016789001 - Айдос Мамедов
-- 77016789002 - Мария Ким
-- 77016789003 - Ерболат Сейтов
-- 77016789004 - Анна Петрова
-- 77016789005 - Нурбол Алдабергенов

-- ============================================
-- barbos_almaty_esentai (4 барбера)
-- ============================================

-- Данияр Абдрахманов (77012345111)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'barbos_almaty_esentai', data.client_phone, '77012345111', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    -- Эта неделя (5 записей)
    ('77016789001', 'created', '2026-01-13 10:00:00+06'::timestamptz, '2026-01-13 10:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789002', 'completed', '2026-01-14 14:00:00+06'::timestamptz, '2026-01-14 14:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789003', 'created', '2026-01-16 11:00:00+06'::timestamptz, '2026-01-16 11:30:00+06'::timestamptz, 'Стрижка бороды'),
    ('77016789004', 'pending', '2026-01-18 15:00:00+06'::timestamptz, '2026-01-18 15:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789005', 'created', '2026-01-19 12:00:00+06'::timestamptz, '2026-01-19 12:50:00+06'::timestamptz, 'Модельная стрижка'),
    -- Следующая неделя (3 записи)
    ('77016789001', 'pending', '2026-01-20 10:00:00+06'::timestamptz, '2026-01-20 10:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789003', 'pending', '2026-01-22 16:00:00+06'::timestamptz, '2026-01-22 16:30:00+06'::timestamptz, 'Стрижка бороды'),
    ('77016789005', 'pending', '2026-01-24 14:00:00+06'::timestamptz, '2026-01-24 14:30:00+06'::timestamptz, 'Детская стрижка')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'barbos_almaty_esentai' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77012345111' AND ms.service_id = s.id;

-- Максим Ковалев (77012345112)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'barbos_almaty_esentai', data.client_phone, '77012345112', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789002', 'completed', '2026-01-13 11:00:00+06'::timestamptz, '2026-01-13 11:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789003', 'created', '2026-01-14 16:00:00+06'::timestamptz, '2026-01-14 17:00:00+06'::timestamptz, 'Королевское бритье'),
    ('77016789001', 'canceled', '2026-01-15 10:00:00+06'::timestamptz, '2026-01-15 10:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789004', 'created', '2026-01-18 11:00:00+06'::timestamptz, '2026-01-18 11:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789005', 'pending', '2026-01-19 14:00:00+06'::timestamptz, '2026-01-19 14:40:00+06'::timestamptz, 'Моделирование бороды'),
    ('77016789001', 'pending', '2026-01-21 10:00:00+06'::timestamptz, '2026-01-21 11:00:00+06'::timestamptz, 'Королевское бритье'),
    ('77016789002', 'pending', '2026-01-23 15:00:00+06'::timestamptz, '2026-01-23 15:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789004', 'pending', '2026-01-25 12:00:00+06'::timestamptz, '2026-01-25 12:40:00+06'::timestamptz, 'Моделирование бороды')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'barbos_almaty_esentai' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77012345112' AND ms.service_id = s.id;

-- Тимур Жолдасов (77012345113)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'barbos_almaty_esentai', data.client_phone, '77012345113', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789001', 'completed', '2026-01-13 14:00:00+06'::timestamptz, '2026-01-13 14:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789004', 'completed', '2026-01-14 10:00:00+06'::timestamptz, '2026-01-14 10:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789002', 'created', '2026-01-16 15:00:00+06'::timestamptz, '2026-01-16 15:45:00+06'::timestamptz, 'Окрашивание бороды'),
    ('77016789003', 'pending', '2026-01-18 13:00:00+06'::timestamptz, '2026-01-18 13:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789005', 'created', '2026-01-19 10:00:00+06'::timestamptz, '2026-01-19 10:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789002', 'pending', '2026-01-20 14:00:00+06'::timestamptz, '2026-01-20 14:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789003', 'pending', '2026-01-22 11:00:00+06'::timestamptz, '2026-01-22 11:45:00+06'::timestamptz, 'Окрашивание бороды'),
    ('77016789001', 'pending', '2026-01-24 16:00:00+06'::timestamptz, '2026-01-24 16:30:00+06'::timestamptz, 'Детская стрижка')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'barbos_almaty_esentai' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77012345113' AND ms.service_id = s.id;

-- Александр Сидоров (77012345114)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'barbos_almaty_esentai', data.client_phone, '77012345114', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789005', 'completed', '2026-01-13 16:00:00+06'::timestamptz, '2026-01-13 17:00:00+06'::timestamptz, 'Королевское бритье'),
    ('77016789001', 'completed', '2026-01-14 11:00:00+06'::timestamptz, '2026-01-14 12:00:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789003', 'canceled', '2026-01-15 14:00:00+06'::timestamptz, '2026-01-15 14:30:00+06'::timestamptz, 'Стрижка бороды'),
    ('77016789002', 'created', '2026-01-18 17:00:00+06'::timestamptz, '2026-01-18 18:00:00+06'::timestamptz, 'Королевское бритье'),
    ('77016789004', 'pending', '2026-01-19 15:00:00+06'::timestamptz, '2026-01-19 16:00:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789005', 'pending', '2026-01-21 12:00:00+06'::timestamptz, '2026-01-21 12:30:00+06'::timestamptz, 'Стрижка бороды'),
    ('77016789001', 'pending', '2026-01-23 10:00:00+06'::timestamptz, '2026-01-23 11:00:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789003', 'pending', '2026-01-25 14:00:00+06'::timestamptz, '2026-01-25 15:00:00+06'::timestamptz, 'Королевское бритье')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'barbos_almaty_esentai' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77012345114' AND ms.service_id = s.id;

-- ============================================
-- barbos_almaty_mega (3 барбера)
-- ============================================

-- Ержан Омаров (77012345211)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'barbos_almaty_mega', data.client_phone, '77012345211', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789001', 'completed', '2026-01-13 10:00:00+06'::timestamptz, '2026-01-13 10:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789002', 'created', '2026-01-15 14:00:00+06'::timestamptz, '2026-01-15 15:00:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789003', 'completed', '2026-01-16 11:00:00+06'::timestamptz, '2026-01-16 11:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789004', 'pending', '2026-01-18 15:00:00+06'::timestamptz, '2026-01-18 15:40:00+06'::timestamptz, 'Классическое бритье'),
    ('77016789005', 'created', '2026-01-19 12:00:00+06'::timestamptz, '2026-01-19 12:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789001', 'pending', '2026-01-20 10:00:00+06'::timestamptz, '2026-01-20 10:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789003', 'pending', '2026-01-22 16:00:00+06'::timestamptz, '2026-01-22 16:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789005', 'pending', '2026-01-24 14:00:00+06'::timestamptz, '2026-01-24 14:40:00+06'::timestamptz, 'Мужская классическая стрижка')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'barbos_almaty_mega' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77012345211' AND ms.service_id = s.id;

-- Руслан Петров (77012345212)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'barbos_almaty_mega', data.client_phone, '77012345212', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789002', 'completed', '2026-01-13 11:00:00+06'::timestamptz, '2026-01-13 11:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789003', 'created', '2026-01-14 16:00:00+06'::timestamptz, '2026-01-14 16:30:00+06'::timestamptz, 'Стрижка бороды'),
    ('77016789001', 'canceled', '2026-01-15 10:00:00+06'::timestamptz, '2026-01-15 10:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789004', 'created', '2026-01-18 11:00:00+06'::timestamptz, '2026-01-18 11:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789005', 'pending', '2026-01-19 14:00:00+06'::timestamptz, '2026-01-19 14:30:00+06'::timestamptz, 'Стрижка бороды'),
    ('77016789001', 'pending', '2026-01-21 10:00:00+06'::timestamptz, '2026-01-21 10:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789002', 'pending', '2026-01-23 15:00:00+06'::timestamptz, '2026-01-23 15:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789004', 'pending', '2026-01-25 12:00:00+06'::timestamptz, '2026-01-25 12:30:00+06'::timestamptz, 'Стрижка бороды')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'barbos_almaty_mega' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77012345212' AND ms.service_id = s.id;

-- Серик Ахметов (77012345213)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'barbos_almaty_mega', data.client_phone, '77012345213', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789001', 'completed', '2026-01-13 14:00:00+06'::timestamptz, '2026-01-13 14:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789004', 'completed', '2026-01-14 10:00:00+06'::timestamptz, '2026-01-14 10:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789002', 'created', '2026-01-16 15:00:00+06'::timestamptz, '2026-01-16 15:40:00+06'::timestamptz, 'Моделирование бороды'),
    ('77016789003', 'pending', '2026-01-18 13:00:00+06'::timestamptz, '2026-01-18 13:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789005', 'created', '2026-01-19 10:00:00+06'::timestamptz, '2026-01-19 10:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789002', 'pending', '2026-01-20 14:00:00+06'::timestamptz, '2026-01-20 14:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789003', 'pending', '2026-01-22 11:00:00+06'::timestamptz, '2026-01-22 11:40:00+06'::timestamptz, 'Моделирование бороды'),
    ('77016789001', 'pending', '2026-01-24 16:00:00+06'::timestamptz, '2026-01-24 16:30:00+06'::timestamptz, 'Детская стрижка')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'barbos_almaty_mega' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77012345213' AND ms.service_id = s.id;

-- ============================================
-- barbos_astana_mega (3 барбера)
-- ============================================

-- Асхат Досмагамбетов (77012345311)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'barbos_astana_mega', data.client_phone, '77012345311', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789001', 'completed', '2026-01-13 10:00:00+06'::timestamptz, '2026-01-13 10:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789002', 'created', '2026-01-15 14:00:00+06'::timestamptz, '2026-01-15 15:00:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789003', 'completed', '2026-01-16 11:00:00+06'::timestamptz, '2026-01-16 11:40:00+06'::timestamptz, 'Классическое бритье'),
    ('77016789004', 'pending', '2026-01-18 15:00:00+06'::timestamptz, '2026-01-18 15:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789005', 'created', '2026-01-19 12:00:00+06'::timestamptz, '2026-01-19 13:00:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789001', 'pending', '2026-01-20 10:00:00+06'::timestamptz, '2026-01-20 10:40:00+06'::timestamptz, 'Классическое бритье'),
    ('77016789003', 'pending', '2026-01-22 16:00:00+06'::timestamptz, '2026-01-22 17:00:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789005', 'pending', '2026-01-24 14:00:00+06'::timestamptz, '2026-01-24 14:40:00+06'::timestamptz, 'Мужская классическая стрижка')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'barbos_astana_mega' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77012345311' AND ms.service_id = s.id;

-- Игорь Васильев (77012345312)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'barbos_astana_mega', data.client_phone, '77012345312', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789002', 'completed', '2026-01-13 11:00:00+06'::timestamptz, '2026-01-13 12:00:00+06'::timestamptz, 'Королевское бритье'),
    ('77016789003', 'created', '2026-01-14 16:00:00+06'::timestamptz, '2026-01-14 16:30:00+06'::timestamptz, 'Стрижка бороды'),
    ('77016789001', 'canceled', '2026-01-15 10:00:00+06'::timestamptz, '2026-01-15 10:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789004', 'created', '2026-01-18 11:00:00+06'::timestamptz, '2026-01-18 11:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789005', 'pending', '2026-01-19 14:00:00+06'::timestamptz, '2026-01-19 15:00:00+06'::timestamptz, 'Королевское бритье'),
    ('77016789001', 'pending', '2026-01-21 10:00:00+06'::timestamptz, '2026-01-21 10:30:00+06'::timestamptz, 'Стрижка бороды'),
    ('77016789002', 'pending', '2026-01-23 15:00:00+06'::timestamptz, '2026-01-23 16:00:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789004', 'pending', '2026-01-25 12:00:00+06'::timestamptz, '2026-01-25 13:00:00+06'::timestamptz, 'Королевское бритье')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'barbos_astana_mega' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77012345312' AND ms.service_id = s.id;

-- Ринат Сулейменов (77012345313)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'barbos_astana_mega', data.client_phone, '77012345313', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789001', 'completed', '2026-01-13 14:00:00+06'::timestamptz, '2026-01-13 15:00:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789004', 'completed', '2026-01-14 10:00:00+06'::timestamptz, '2026-01-14 10:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789002', 'created', '2026-01-16 15:00:00+06'::timestamptz, '2026-01-16 15:40:00+06'::timestamptz, 'Моделирование бороды'),
    ('77016789003', 'pending', '2026-01-18 13:00:00+06'::timestamptz, '2026-01-18 14:00:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789005', 'created', '2026-01-19 10:00:00+06'::timestamptz, '2026-01-19 10:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789002', 'pending', '2026-01-20 14:00:00+06'::timestamptz, '2026-01-20 15:00:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789003', 'pending', '2026-01-22 11:00:00+06'::timestamptz, '2026-01-22 11:40:00+06'::timestamptz, 'Моделирование бороды'),
    ('77016789001', 'pending', '2026-01-24 16:00:00+06'::timestamptz, '2026-01-24 16:30:00+06'::timestamptz, 'Детская стрижка')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'barbos_astana_mega' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77012345313' AND ms.service_id = s.id;

-- ============================================
-- chic_almaty_furmanov (5 мастеров)
-- ============================================

-- Асем Нурланова (77013456111) - стрижки и укладки
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'chic_almaty_furmanov', data.client_phone, '77013456111', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789002', 'completed', '2026-01-13 09:00:00+06'::timestamptz, '2026-01-13 10:00:00+06'::timestamptz, 'Женская стрижка'),
    ('77016789004', 'created', '2026-01-15 14:00:00+06'::timestamptz, '2026-01-15 14:40:00+06'::timestamptz, 'Укладка феном'),
    ('77016789002', 'completed', '2026-01-16 11:00:00+06'::timestamptz, '2026-01-16 12:30:00+06'::timestamptz, 'Вечерняя укладка'),
    ('77016789004', 'pending', '2026-01-18 15:00:00+06'::timestamptz, '2026-01-18 16:00:00+06'::timestamptz, 'Женская стрижка'),
    ('77016789002', 'created', '2026-01-19 12:00:00+06'::timestamptz, '2026-01-19 12:40:00+06'::timestamptz, 'Укладка феном'),
    ('77016789004', 'pending', '2026-01-20 10:00:00+06'::timestamptz, '2026-01-20 11:00:00+06'::timestamptz, 'Женская стрижка'),
    ('77016789002', 'pending', '2026-01-22 16:00:00+06'::timestamptz, '2026-01-22 17:30:00+06'::timestamptz, 'Вечерняя укладка'),
    ('77016789004', 'pending', '2026-01-24 14:00:00+06'::timestamptz, '2026-01-24 14:40:00+06'::timestamptz, 'Укладка феном')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'chic_almaty_furmanov' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77013456111' AND ms.service_id = s.id;

-- Дина Жумабаева (77013456112) - окрашивание
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'chic_almaty_furmanov', data.client_phone, '77013456112', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789002', 'completed', '2026-01-13 10:00:00+06'::timestamptz, '2026-01-13 12:00:00+06'::timestamptz, 'Полное окрашивание'),
    ('77016789004', 'created', '2026-01-14 14:00:00+06'::timestamptz, '2026-01-14 16:30:00+06'::timestamptz, 'Мелирование'),
    ('77016789002', 'completed', '2026-01-15 11:00:00+06'::timestamptz, '2026-01-15 14:00:00+06'::timestamptz, 'Балаяж'),
    ('77016789004', 'pending', '2026-01-18 10:00:00+06'::timestamptz, '2026-01-18 12:00:00+06'::timestamptz, 'Полное окрашивание'),
    ('77016789002', 'created', '2026-01-19 14:00:00+06'::timestamptz, '2026-01-19 17:00:00+06'::timestamptz, 'Омбре'),
    ('77016789004', 'pending', '2026-01-21 10:00:00+06'::timestamptz, '2026-01-21 12:30:00+06'::timestamptz, 'Мелирование'),
    ('77016789002', 'pending', '2026-01-23 11:00:00+06'::timestamptz, '2026-01-23 14:00:00+06'::timestamptz, 'Балаяж'),
    ('77016789004', 'pending', '2026-01-25 10:00:00+06'::timestamptz, '2026-01-25 12:00:00+06'::timestamptz, 'Полное окрашивание')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'chic_almaty_furmanov' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77013456112' AND ms.service_id = s.id;

-- Елена Смирнова (77013456113) - маникюр
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'chic_almaty_furmanov', data.client_phone, '77013456113', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789002', 'completed', '2026-01-13 09:00:00+06'::timestamptz, '2026-01-13 10:00:00+06'::timestamptz, 'Классический маникюр'),
    ('77016789004', 'created', '2026-01-14 11:00:00+06'::timestamptz, '2026-01-14 12:00:00+06'::timestamptz, 'Аппаратный маникюр'),
    ('77016789002', 'completed', '2026-01-16 15:00:00+06'::timestamptz, '2026-01-16 16:30:00+06'::timestamptz, 'Маникюр с покрытием'),
    ('77016789004', 'pending', '2026-01-18 14:00:00+06'::timestamptz, '2026-01-18 15:00:00+06'::timestamptz, 'Классический маникюр'),
    ('77016789002', 'created', '2026-01-19 10:00:00+06'::timestamptz, '2026-01-19 11:30:00+06'::timestamptz, 'Маникюр с покрытием'),
    ('77016789004', 'pending', '2026-01-20 15:00:00+06'::timestamptz, '2026-01-20 16:00:00+06'::timestamptz, 'Аппаратный маникюр'),
    ('77016789002', 'pending', '2026-01-22 10:00:00+06'::timestamptz, '2026-01-22 11:30:00+06'::timestamptz, 'Маникюр с покрытием'),
    ('77016789004', 'pending', '2026-01-24 11:00:00+06'::timestamptz, '2026-01-24 12:00:00+06'::timestamptz, 'Классический маникюр')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'chic_almaty_furmanov' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77013456113' AND ms.service_id = s.id;

-- Жанна Бекмуратова (77013456114) - педикюр
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'chic_almaty_furmanov', data.client_phone, '77013456114', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789002', 'completed', '2026-01-13 10:00:00+06'::timestamptz, '2026-01-13 11:00:00+06'::timestamptz, 'Классический педикюр'),
    ('77016789004', 'created', '2026-01-14 14:00:00+06'::timestamptz, '2026-01-14 15:00:00+06'::timestamptz, 'Аппаратный педикюр'),
    ('77016789002', 'completed', '2026-01-15 11:00:00+06'::timestamptz, '2026-01-15 12:30:00+06'::timestamptz, 'SPA-педикюр'),
    ('77016789004', 'pending', '2026-01-18 16:00:00+06'::timestamptz, '2026-01-18 17:00:00+06'::timestamptz, 'Классический педикюр'),
    ('77016789002', 'created', '2026-01-19 15:00:00+06'::timestamptz, '2026-01-19 16:30:00+06'::timestamptz, 'SPA-педикюр'),
    ('77016789004', 'pending', '2026-01-21 14:00:00+06'::timestamptz, '2026-01-21 15:00:00+06'::timestamptz, 'Аппаратный педикюр'),
    ('77016789002', 'pending', '2026-01-23 10:00:00+06'::timestamptz, '2026-01-23 11:30:00+06'::timestamptz, 'SPA-педикюр'),
    ('77016789004', 'pending', '2026-01-25 16:00:00+06'::timestamptz, '2026-01-25 17:00:00+06'::timestamptz, 'Классический педикюр')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'chic_almaty_furmanov' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77013456114' AND ms.service_id = s.id;

-- Индира Касымова (77013456115) - косметология
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'chic_almaty_furmanov', data.client_phone, '77013456115', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789002', 'completed', '2026-01-13 11:00:00+06'::timestamptz, '2026-01-13 12:00:00+06'::timestamptz, 'Чистка лица'),
    ('77016789004', 'created', '2026-01-14 15:00:00+06'::timestamptz, '2026-01-14 16:00:00+06'::timestamptz, 'Уходовая процедура'),
    ('77016789002', 'completed', '2026-01-16 10:00:00+06'::timestamptz, '2026-01-16 11:00:00+06'::timestamptz, 'Чистка лица'),
    ('77016789004', 'pending', '2026-01-18 12:00:00+06'::timestamptz, '2026-01-18 13:00:00+06'::timestamptz, 'Уходовая процедура'),
    ('77016789002', 'created', '2026-01-19 11:00:00+06'::timestamptz, '2026-01-19 12:00:00+06'::timestamptz, 'Чистка лица'),
    ('77016789004', 'pending', '2026-01-20 16:00:00+06'::timestamptz, '2026-01-20 17:00:00+06'::timestamptz, 'Уходовая процедура'),
    ('77016789002', 'pending', '2026-01-22 14:00:00+06'::timestamptz, '2026-01-22 15:00:00+06'::timestamptz, 'Чистка лица'),
    ('77016789004', 'pending', '2026-01-24 10:00:00+06'::timestamptz, '2026-01-24 11:00:00+06'::timestamptz, 'Уходовая процедура')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'chic_almaty_furmanov' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77013456115' AND ms.service_id = s.id;

-- ============================================
-- chic_almaty_dostyk (4 мастера)
-- ============================================

-- Камила Оспанова (77013456211)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'chic_almaty_dostyk', data.client_phone, '77013456211', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789002', 'completed', '2026-01-13 09:00:00+06'::timestamptz, '2026-01-13 10:00:00+06'::timestamptz, 'Женская стрижка'),
    ('77016789004', 'created', '2026-01-15 14:00:00+06'::timestamptz, '2026-01-15 14:40:00+06'::timestamptz, 'Укладка феном'),
    ('77016789002', 'completed', '2026-01-16 11:00:00+06'::timestamptz, '2026-01-16 12:30:00+06'::timestamptz, 'Вечерняя укладка'),
    ('77016789004', 'pending', '2026-01-18 15:00:00+06'::timestamptz, '2026-01-18 16:00:00+06'::timestamptz, 'Женская стрижка'),
    ('77016789002', 'created', '2026-01-19 12:00:00+06'::timestamptz, '2026-01-19 12:40:00+06'::timestamptz, 'Укладка феном'),
    ('77016789004', 'pending', '2026-01-20 10:00:00+06'::timestamptz, '2026-01-20 11:00:00+06'::timestamptz, 'Женская стрижка'),
    ('77016789002', 'pending', '2026-01-22 16:00:00+06'::timestamptz, '2026-01-22 17:30:00+06'::timestamptz, 'Вечерняя укладка'),
    ('77016789004', 'pending', '2026-01-24 14:00:00+06'::timestamptz, '2026-01-24 14:40:00+06'::timestamptz, 'Укладка феном')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'chic_almaty_dostyk' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77013456211' AND ms.service_id = s.id;

-- Лаура Мухамбетова (77013456212)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'chic_almaty_dostyk', data.client_phone, '77013456212', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789002', 'completed', '2026-01-13 10:00:00+06'::timestamptz, '2026-01-13 12:00:00+06'::timestamptz, 'Полное окрашивание'),
    ('77016789004', 'created', '2026-01-15 14:00:00+06'::timestamptz, '2026-01-15 17:00:00+06'::timestamptz, 'Балаяж'),
    ('77016789002', 'completed', '2026-01-16 11:00:00+06'::timestamptz, '2026-01-16 13:00:00+06'::timestamptz, 'Полное окрашивание'),
    ('77016789004', 'pending', '2026-01-18 10:00:00+06'::timestamptz, '2026-01-18 13:00:00+06'::timestamptz, 'Балаяж'),
    ('77016789002', 'created', '2026-01-19 14:00:00+06'::timestamptz, '2026-01-19 16:00:00+06'::timestamptz, 'Полное окрашивание'),
    ('77016789004', 'pending', '2026-01-21 10:00:00+06'::timestamptz, '2026-01-21 13:00:00+06'::timestamptz, 'Балаяж'),
    ('77016789002', 'pending', '2026-01-23 11:00:00+06'::timestamptz, '2026-01-23 13:00:00+06'::timestamptz, 'Полное окрашивание'),
    ('77016789004', 'pending', '2026-01-25 10:00:00+06'::timestamptz, '2026-01-25 13:00:00+06'::timestamptz, 'Балаяж')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'chic_almaty_dostyk' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77013456212' AND ms.service_id = s.id;

-- Марина Иванова (77013456213)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'chic_almaty_dostyk', data.client_phone, '77013456213', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789002', 'completed', '2026-01-13 09:00:00+06'::timestamptz, '2026-01-13 10:00:00+06'::timestamptz, 'Аппаратный маникюр'),
    ('77016789004', 'created', '2026-01-14 11:00:00+06'::timestamptz, '2026-01-14 12:30:00+06'::timestamptz, 'Маникюр с покрытием'),
    ('77016789002', 'completed', '2026-01-16 15:00:00+06'::timestamptz, '2026-01-16 16:00:00+06'::timestamptz, 'Аппаратный маникюр'),
    ('77016789004', 'pending', '2026-01-18 14:00:00+06'::timestamptz, '2026-01-18 15:30:00+06'::timestamptz, 'Маникюр с покрытием'),
    ('77016789002', 'created', '2026-01-19 10:00:00+06'::timestamptz, '2026-01-19 11:00:00+06'::timestamptz, 'Аппаратный маникюр'),
    ('77016789004', 'pending', '2026-01-20 15:00:00+06'::timestamptz, '2026-01-20 16:30:00+06'::timestamptz, 'Маникюр с покрытием'),
    ('77016789002', 'pending', '2026-01-22 10:00:00+06'::timestamptz, '2026-01-22 11:00:00+06'::timestamptz, 'Аппаратный маникюр'),
    ('77016789004', 'pending', '2026-01-24 11:00:00+06'::timestamptz, '2026-01-24 12:30:00+06'::timestamptz, 'Маникюр с покрытием')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'chic_almaty_dostyk' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77013456213' AND ms.service_id = s.id;

-- Назым Қаржаубаева (77013456214)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'chic_almaty_dostyk', data.client_phone, '77013456214', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789002', 'completed', '2026-01-13 10:00:00+06'::timestamptz, '2026-01-13 11:00:00+06'::timestamptz, 'Аппаратный педикюр'),
    ('77016789004', 'created', '2026-01-14 14:00:00+06'::timestamptz, '2026-01-14 15:00:00+06'::timestamptz, 'Чистка лица'),
    ('77016789002', 'completed', '2026-01-15 11:00:00+06'::timestamptz, '2026-01-15 12:00:00+06'::timestamptz, 'Аппаратный педикюр'),
    ('77016789004', 'pending', '2026-01-18 16:00:00+06'::timestamptz, '2026-01-18 17:00:00+06'::timestamptz, 'Чистка лица'),
    ('77016789002', 'created', '2026-01-19 15:00:00+06'::timestamptz, '2026-01-19 16:00:00+06'::timestamptz, 'Аппаратный педикюр'),
    ('77016789004', 'pending', '2026-01-21 14:00:00+06'::timestamptz, '2026-01-21 15:00:00+06'::timestamptz, 'Чистка лица'),
    ('77016789002', 'pending', '2026-01-23 10:00:00+06'::timestamptz, '2026-01-23 11:00:00+06'::timestamptz, 'Аппаратный педикюр'),
    ('77016789004', 'pending', '2026-01-25 16:00:00+06'::timestamptz, '2026-01-25 17:00:00+06'::timestamptz, 'Чистка лица')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'chic_almaty_dostyk' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77013456214' AND ms.service_id = s.id;

-- ============================================
-- brocode_almaty_abay (3 барбера)
-- ============================================

-- Азамат Досанов (77014567111)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'brocode_almaty_abay', data.client_phone, '77014567111', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789001', 'completed', '2026-01-13 10:00:00+06'::timestamptz, '2026-01-13 10:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789003', 'created', '2026-01-15 14:00:00+06'::timestamptz, '2026-01-15 14:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789005', 'completed', '2026-01-16 11:00:00+06'::timestamptz, '2026-01-16 11:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789001', 'pending', '2026-01-18 15:00:00+06'::timestamptz, '2026-01-18 15:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789003', 'created', '2026-01-19 12:00:00+06'::timestamptz, '2026-01-19 12:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789005', 'pending', '2026-01-20 10:00:00+06'::timestamptz, '2026-01-20 10:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789001', 'pending', '2026-01-22 16:00:00+06'::timestamptz, '2026-01-22 16:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789003', 'pending', '2026-01-24 14:00:00+06'::timestamptz, '2026-01-24 14:50:00+06'::timestamptz, 'Модельная стрижка')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'brocode_almaty_abay' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77014567111' AND ms.service_id = s.id;

-- Батыр Есимов (77014567112)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'brocode_almaty_abay', data.client_phone, '77014567112', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789003', 'completed', '2026-01-13 11:00:00+06'::timestamptz, '2026-01-13 11:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789001', 'created', '2026-01-14 16:00:00+06'::timestamptz, '2026-01-14 16:30:00+06'::timestamptz, 'Стрижка бороды'),
    ('77016789005', 'canceled', '2026-01-15 10:00:00+06'::timestamptz, '2026-01-15 10:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789003', 'created', '2026-01-18 11:00:00+06'::timestamptz, '2026-01-18 11:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789001', 'pending', '2026-01-19 14:00:00+06'::timestamptz, '2026-01-19 14:30:00+06'::timestamptz, 'Стрижка бороды'),
    ('77016789005', 'pending', '2026-01-21 10:00:00+06'::timestamptz, '2026-01-21 10:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789003', 'pending', '2026-01-23 15:00:00+06'::timestamptz, '2026-01-23 15:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789001', 'pending', '2026-01-25 12:00:00+06'::timestamptz, '2026-01-25 12:30:00+06'::timestamptz, 'Стрижка бороды')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'brocode_almaty_abay' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77014567112' AND ms.service_id = s.id;

-- Владислав Морозов (77014567113)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'brocode_almaty_abay', data.client_phone, '77014567113', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789005', 'completed', '2026-01-13 14:00:00+06'::timestamptz, '2026-01-13 14:40:00+06'::timestamptz, 'Классическое бритье'),
    ('77016789001', 'completed', '2026-01-14 10:00:00+06'::timestamptz, '2026-01-14 10:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789003', 'created', '2026-01-16 15:00:00+06'::timestamptz, '2026-01-16 15:20:00+06'::timestamptz, 'Укладка волос'),
    ('77016789005', 'pending', '2026-01-18 13:00:00+06'::timestamptz, '2026-01-18 13:40:00+06'::timestamptz, 'Классическое бритье'),
    ('77016789001', 'created', '2026-01-19 10:00:00+06'::timestamptz, '2026-01-19 10:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789003', 'pending', '2026-01-20 14:00:00+06'::timestamptz, '2026-01-20 14:40:00+06'::timestamptz, 'Классическое бритье'),
    ('77016789005', 'pending', '2026-01-22 11:00:00+06'::timestamptz, '2026-01-22 11:20:00+06'::timestamptz, 'Укладка волос'),
    ('77016789001', 'pending', '2026-01-24 16:00:00+06'::timestamptz, '2026-01-24 16:30:00+06'::timestamptz, 'Детская стрижка')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'brocode_almaty_abay' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77014567113' AND ms.service_id = s.id;

-- ============================================
-- brocode_astana_respublika (3 барбера)
-- ============================================

-- Дархан Амангельдин (77014567211)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'brocode_astana_respublika', data.client_phone, '77014567211', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789001', 'completed', '2026-01-13 10:00:00+06'::timestamptz, '2026-01-13 10:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789003', 'created', '2026-01-15 14:00:00+06'::timestamptz, '2026-01-15 14:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789005', 'completed', '2026-01-16 11:00:00+06'::timestamptz, '2026-01-16 11:40:00+06'::timestamptz, 'Классическое бритье'),
    ('77016789001', 'pending', '2026-01-18 15:00:00+06'::timestamptz, '2026-01-18 15:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789003', 'created', '2026-01-19 12:00:00+06'::timestamptz, '2026-01-19 12:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789005', 'pending', '2026-01-20 10:00:00+06'::timestamptz, '2026-01-20 10:40:00+06'::timestamptz, 'Классическое бритье'),
    ('77016789001', 'pending', '2026-01-22 16:00:00+06'::timestamptz, '2026-01-22 16:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789003', 'pending', '2026-01-24 14:00:00+06'::timestamptz, '2026-01-24 14:40:00+06'::timestamptz, 'Мужская классическая стрижка')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'brocode_astana_respublika' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77014567211' AND ms.service_id = s.id;

-- Евгений Соколов (77014567212)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'brocode_astana_respublika', data.client_phone, '77014567212', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789003', 'completed', '2026-01-13 11:00:00+06'::timestamptz, '2026-01-13 11:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789001', 'created', '2026-01-14 16:00:00+06'::timestamptz, '2026-01-14 16:30:00+06'::timestamptz, 'Стрижка бороды'),
    ('77016789005', 'canceled', '2026-01-15 10:00:00+06'::timestamptz, '2026-01-15 10:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789003', 'created', '2026-01-18 11:00:00+06'::timestamptz, '2026-01-18 11:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789001', 'pending', '2026-01-19 14:00:00+06'::timestamptz, '2026-01-19 14:30:00+06'::timestamptz, 'Стрижка бороды'),
    ('77016789005', 'pending', '2026-01-21 10:00:00+06'::timestamptz, '2026-01-21 10:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789003', 'pending', '2026-01-23 15:00:00+06'::timestamptz, '2026-01-23 15:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789001', 'pending', '2026-01-25 12:00:00+06'::timestamptz, '2026-01-25 12:30:00+06'::timestamptz, 'Стрижка бороды')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'brocode_astana_respublika' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77014567212' AND ms.service_id = s.id;

-- Жасулан Қалиев (77014567213)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'brocode_astana_respublika', data.client_phone, '77014567213', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789005', 'completed', '2026-01-13 14:00:00+06'::timestamptz, '2026-01-13 14:40:00+06'::timestamptz, 'Классическое бритье'),
    ('77016789001', 'completed', '2026-01-14 10:00:00+06'::timestamptz, '2026-01-14 10:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789003', 'created', '2026-01-16 15:00:00+06'::timestamptz, '2026-01-16 15:30:00+06'::timestamptz, 'Стрижка бороды'),
    ('77016789005', 'pending', '2026-01-18 13:00:00+06'::timestamptz, '2026-01-18 13:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789001', 'created', '2026-01-19 10:00:00+06'::timestamptz, '2026-01-19 10:40:00+06'::timestamptz, 'Классическое бритье'),
    ('77016789003', 'pending', '2026-01-20 14:00:00+06'::timestamptz, '2026-01-20 14:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789005', 'pending', '2026-01-22 11:00:00+06'::timestamptz, '2026-01-22 11:30:00+06'::timestamptz, 'Стрижка бороды'),
    ('77016789001', 'pending', '2026-01-24 16:00:00+06'::timestamptz, '2026-01-24 16:50:00+06'::timestamptz, 'Модельная стрижка')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'brocode_astana_respublika' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77014567213' AND ms.service_id = s.id;

-- ============================================
-- brocode_shymkent_tauke (3 барбера)
-- ============================================

-- Кайрат Төлеген (77014567311)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'brocode_shymkent_tauke', data.client_phone, '77014567311', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789001', 'completed', '2026-01-13 10:00:00+06'::timestamptz, '2026-01-13 10:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789003', 'created', '2026-01-15 14:00:00+06'::timestamptz, '2026-01-15 14:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789005', 'completed', '2026-01-16 11:00:00+06'::timestamptz, '2026-01-16 11:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789001', 'pending', '2026-01-18 15:00:00+06'::timestamptz, '2026-01-18 15:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789003', 'created', '2026-01-19 12:00:00+06'::timestamptz, '2026-01-19 12:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789005', 'pending', '2026-01-20 10:00:00+06'::timestamptz, '2026-01-20 10:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789001', 'pending', '2026-01-22 16:00:00+06'::timestamptz, '2026-01-22 16:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789003', 'pending', '2026-01-24 14:00:00+06'::timestamptz, '2026-01-24 14:50:00+06'::timestamptz, 'Модельная стрижка')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'brocode_shymkent_tauke' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77014567311' AND ms.service_id = s.id;

-- Марат Ибрагимов (77014567312)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'brocode_shymkent_tauke', data.client_phone, '77014567312', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789003', 'completed', '2026-01-13 11:00:00+06'::timestamptz, '2026-01-13 11:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789001', 'created', '2026-01-14 16:00:00+06'::timestamptz, '2026-01-14 16:30:00+06'::timestamptz, 'Стрижка бороды'),
    ('77016789005', 'canceled', '2026-01-15 10:00:00+06'::timestamptz, '2026-01-15 10:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789003', 'created', '2026-01-18 11:00:00+06'::timestamptz, '2026-01-18 11:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789001', 'pending', '2026-01-19 14:00:00+06'::timestamptz, '2026-01-19 14:30:00+06'::timestamptz, 'Стрижка бороды'),
    ('77016789005', 'pending', '2026-01-21 10:00:00+06'::timestamptz, '2026-01-21 10:40:00+06'::timestamptz, 'Мужская классическая стрижка'),
    ('77016789003', 'pending', '2026-01-23 15:00:00+06'::timestamptz, '2026-01-23 15:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789001', 'pending', '2026-01-25 12:00:00+06'::timestamptz, '2026-01-25 12:30:00+06'::timestamptz, 'Стрижка бороды')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'brocode_shymkent_tauke' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77014567312' AND ms.service_id = s.id;

-- Нуржан Сапаров (77014567313)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'brocode_shymkent_tauke', data.client_phone, '77014567313', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789005', 'completed', '2026-01-13 14:00:00+06'::timestamptz, '2026-01-13 14:40:00+06'::timestamptz, 'Классическое бритье'),
    ('77016789001', 'completed', '2026-01-14 10:00:00+06'::timestamptz, '2026-01-14 10:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789003', 'created', '2026-01-16 15:00:00+06'::timestamptz, '2026-01-16 15:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789005', 'pending', '2026-01-18 13:00:00+06'::timestamptz, '2026-01-18 13:40:00+06'::timestamptz, 'Классическое бритье'),
    ('77016789001', 'created', '2026-01-19 10:00:00+06'::timestamptz, '2026-01-19 10:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789003', 'pending', '2026-01-20 14:00:00+06'::timestamptz, '2026-01-20 14:50:00+06'::timestamptz, 'Модельная стрижка'),
    ('77016789005', 'pending', '2026-01-22 11:00:00+06'::timestamptz, '2026-01-22 11:40:00+06'::timestamptz, 'Классическое бритье'),
    ('77016789001', 'pending', '2026-01-24 16:00:00+06'::timestamptz, '2026-01-24 16:30:00+06'::timestamptz, 'Детская стрижка')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'brocode_shymkent_tauke' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77014567313' AND ms.service_id = s.id;

-- ============================================
-- krasotka_almaty_zhibek (3 мастера)
-- ============================================

-- Алия Ержанова (77015678111)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'krasotka_almaty_zhibek', data.client_phone, '77015678111', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789001', 'completed', '2026-01-13 09:00:00+06'::timestamptz, '2026-01-13 09:40:00+06'::timestamptz, 'Мужская стрижка'),
    ('77016789002', 'created', '2026-01-15 14:00:00+06'::timestamptz, '2026-01-15 15:00:00+06'::timestamptz, 'Женская стрижка'),
    ('77016789004', 'completed', '2026-01-16 11:00:00+06'::timestamptz, '2026-01-16 11:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789001', 'pending', '2026-01-18 15:00:00+06'::timestamptz, '2026-01-18 15:40:00+06'::timestamptz, 'Мужская стрижка'),
    ('77016789002', 'created', '2026-01-19 12:00:00+06'::timestamptz, '2026-01-19 13:00:00+06'::timestamptz, 'Женская стрижка'),
    ('77016789004', 'pending', '2026-01-20 10:00:00+06'::timestamptz, '2026-01-20 10:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789001', 'pending', '2026-01-22 16:00:00+06'::timestamptz, '2026-01-22 16:40:00+06'::timestamptz, 'Мужская стрижка'),
    ('77016789002', 'pending', '2026-01-24 14:00:00+06'::timestamptz, '2026-01-24 15:00:00+06'::timestamptz, 'Женская стрижка')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'krasotka_almaty_zhibek' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77015678111' AND ms.service_id = s.id;

-- Балжан Қадырова (77015678112)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'krasotka_almaty_zhibek', data.client_phone, '77015678112', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789002', 'completed', '2026-01-13 10:00:00+06'::timestamptz, '2026-01-13 10:30:00+06'::timestamptz, 'Укладка волос'),
    ('77016789004', 'created', '2026-01-14 14:00:00+06'::timestamptz, '2026-01-14 16:00:00+06'::timestamptz, 'Полное окрашивание'),
    ('77016789002', 'completed', '2026-01-15 11:00:00+06'::timestamptz, '2026-01-15 12:00:00+06'::timestamptz, 'Женская стрижка'),
    ('77016789004', 'pending', '2026-01-18 16:00:00+06'::timestamptz, '2026-01-18 16:30:00+06'::timestamptz, 'Укладка волос'),
    ('77016789002', 'created', '2026-01-19 14:00:00+06'::timestamptz, '2026-01-19 16:00:00+06'::timestamptz, 'Полное окрашивание'),
    ('77016789004', 'pending', '2026-01-21 10:00:00+06'::timestamptz, '2026-01-21 11:00:00+06'::timestamptz, 'Женская стрижка'),
    ('77016789002', 'pending', '2026-01-23 15:00:00+06'::timestamptz, '2026-01-23 15:30:00+06'::timestamptz, 'Укладка волос'),
    ('77016789004', 'pending', '2026-01-25 12:00:00+06'::timestamptz, '2026-01-25 14:00:00+06'::timestamptz, 'Полное окрашивание')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'krasotka_almaty_zhibek' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77015678112' AND ms.service_id = s.id;

-- Виктория Сергеева (77015678113)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'krasotka_almaty_zhibek', data.client_phone, '77015678113', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789002', 'completed', '2026-01-13 09:00:00+06'::timestamptz, '2026-01-13 10:00:00+06'::timestamptz, 'Классический маникюр'),
    ('77016789004', 'created', '2026-01-14 11:00:00+06'::timestamptz, '2026-01-14 12:00:00+06'::timestamptz, 'Классический педикюр'),
    ('77016789002', 'completed', '2026-01-16 15:00:00+06'::timestamptz, '2026-01-16 16:00:00+06'::timestamptz, 'Классический маникюр'),
    ('77016789004', 'pending', '2026-01-18 14:00:00+06'::timestamptz, '2026-01-18 15:00:00+06'::timestamptz, 'Классический педикюр'),
    ('77016789002', 'created', '2026-01-19 10:00:00+06'::timestamptz, '2026-01-19 11:00:00+06'::timestamptz, 'Классический маникюр'),
    ('77016789004', 'pending', '2026-01-20 15:00:00+06'::timestamptz, '2026-01-20 16:00:00+06'::timestamptz, 'Классический педикюр'),
    ('77016789002', 'pending', '2026-01-22 10:00:00+06'::timestamptz, '2026-01-22 11:00:00+06'::timestamptz, 'Классический маникюр'),
    ('77016789004', 'pending', '2026-01-24 11:00:00+06'::timestamptz, '2026-01-24 12:00:00+06'::timestamptz, 'Классический педикюр')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'krasotka_almaty_zhibek' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77015678113' AND ms.service_id = s.id;

-- ============================================
-- krasotka_almaty_alatau (3 мастера)
-- ============================================

-- Динара Смагулова (77015678211)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'krasotka_almaty_alatau', data.client_phone, '77015678211', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789001', 'completed', '2026-01-13 09:00:00+06'::timestamptz, '2026-01-13 09:40:00+06'::timestamptz, 'Мужская стрижка'),
    ('77016789002', 'created', '2026-01-15 14:00:00+06'::timestamptz, '2026-01-15 15:00:00+06'::timestamptz, 'Женская стрижка'),
    ('77016789004', 'completed', '2026-01-16 11:00:00+06'::timestamptz, '2026-01-16 11:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789001', 'pending', '2026-01-18 15:00:00+06'::timestamptz, '2026-01-18 15:40:00+06'::timestamptz, 'Мужская стрижка'),
    ('77016789002', 'created', '2026-01-19 12:00:00+06'::timestamptz, '2026-01-19 13:00:00+06'::timestamptz, 'Женская стрижка'),
    ('77016789004', 'pending', '2026-01-20 10:00:00+06'::timestamptz, '2026-01-20 10:30:00+06'::timestamptz, 'Детская стрижка'),
    ('77016789001', 'pending', '2026-01-22 16:00:00+06'::timestamptz, '2026-01-22 16:40:00+06'::timestamptz, 'Мужская стрижка'),
    ('77016789002', 'pending', '2026-01-24 14:00:00+06'::timestamptz, '2026-01-24 15:00:00+06'::timestamptz, 'Женская стрижка')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'krasotka_almaty_alatau' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77015678211' AND ms.service_id = s.id;

-- Екатерина Волкова (77015678212)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'krasotka_almaty_alatau', data.client_phone, '77015678212', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789002', 'completed', '2026-01-13 10:00:00+06'::timestamptz, '2026-01-13 10:30:00+06'::timestamptz, 'Укладка волос'),
    ('77016789004', 'created', '2026-01-14 14:00:00+06'::timestamptz, '2026-01-14 16:00:00+06'::timestamptz, 'Полное окрашивание'),
    ('77016789002', 'completed', '2026-01-15 11:00:00+06'::timestamptz, '2026-01-15 12:00:00+06'::timestamptz, 'Женская стрижка'),
    ('77016789004', 'pending', '2026-01-18 16:00:00+06'::timestamptz, '2026-01-18 16:30:00+06'::timestamptz, 'Укладка волос'),
    ('77016789002', 'created', '2026-01-19 14:00:00+06'::timestamptz, '2026-01-19 16:00:00+06'::timestamptz, 'Полное окрашивание'),
    ('77016789004', 'pending', '2026-01-21 10:00:00+06'::timestamptz, '2026-01-21 11:00:00+06'::timestamptz, 'Женская стрижка'),
    ('77016789002', 'pending', '2026-01-23 15:00:00+06'::timestamptz, '2026-01-23 15:30:00+06'::timestamptz, 'Укладка волос'),
    ('77016789004', 'pending', '2026-01-25 12:00:00+06'::timestamptz, '2026-01-25 14:00:00+06'::timestamptz, 'Полное окрашивание')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'krasotka_almaty_alatau' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77015678212' AND ms.service_id = s.id;

-- Жанар Туғанбаева (77015678213)
INSERT INTO bagsies (service_id, point_code, client_phone, master_phone, status, price, start_at, end_at, updated_by)
SELECT s.id, 'krasotka_almaty_alatau', data.client_phone, '77015678213', data.status, ms.price, data.start_at, data.end_at, 'bagsy_seed'
FROM (VALUES
    ('77016789002', 'completed', '2026-01-13 09:00:00+06'::timestamptz, '2026-01-13 10:00:00+06'::timestamptz, 'Аппаратный маникюр'),
    ('77016789004', 'created', '2026-01-14 11:00:00+06'::timestamptz, '2026-01-14 12:00:00+06'::timestamptz, 'Аппаратный педикюр'),
    ('77016789002', 'completed', '2026-01-16 15:00:00+06'::timestamptz, '2026-01-16 16:00:00+06'::timestamptz, 'Аппаратный маникюр'),
    ('77016789004', 'pending', '2026-01-18 14:00:00+06'::timestamptz, '2026-01-18 15:00:00+06'::timestamptz, 'Аппаратный педикюр'),
    ('77016789002', 'created', '2026-01-19 10:00:00+06'::timestamptz, '2026-01-19 11:00:00+06'::timestamptz, 'Аппаратный маникюр'),
    ('77016789004', 'pending', '2026-01-20 15:00:00+06'::timestamptz, '2026-01-20 16:00:00+06'::timestamptz, 'Аппаратный педикюр'),
    ('77016789002', 'pending', '2026-01-22 10:00:00+06'::timestamptz, '2026-01-22 11:00:00+06'::timestamptz, 'Аппаратный маникюр'),
    ('77016789004', 'pending', '2026-01-24 11:00:00+06'::timestamptz, '2026-01-24 12:00:00+06'::timestamptz, 'Аппаратный педикюр')
) AS data(client_phone, status, start_at, end_at, service_name)
JOIN services s ON s.point_code = 'krasotka_almaty_alatau' AND s.name = data.service_name
JOIN master_services ms ON ms.master_phone = '77015678213' AND ms.service_id = s.id;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM bagsies WHERE updated_by = 'bagsy_seed';
-- +goose StatementEnd
