-- +goose Up
-- +goose StatementBegin

-- ============================================
-- ТЕСТОВЫЕ ДАННЫЕ ДЛЯ ПРОВЕРКИ РОЛЕЙ
-- ============================================
--
-- ВАЖНО: Пароль для всех тестовых пользователей: password123
-- Хеш: $2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe
--
-- ============================================

-- Создаем тестовые сети
INSERT INTO networks (code, name, description, updated_by) VALUES
('NET_MANAGER', 'Сеть для NetManager', 'Тестовая сеть для проверки прав NetManager', 'test_seed'),
('NET_SELF_OWNER', 'Сеть для SelfOwner', 'Тестовая сеть для проверки прав SelfOwner', 'test_seed'),
('NET_REGULAR', 'Обычная сеть', 'Тестовая сеть для проверки прав Manager', 'test_seed')
ON CONFLICT (code) DO NOTHING;

-- Создаем категорию точек (если еще не существует)
INSERT INTO point_categories (name, description, updated_by) VALUES
('Барбершоп', 'Барбершоп - мужская парикмахерская', 'test_seed')
ON CONFLICT (name) DO NOTHING;

-- Создаем точки для сети NetManager
INSERT INTO points (code, name, network_code, category_id, city, address, schedule, active, updated_by) VALUES
('POINT_NM_01', 'NetManager Точка 1', 'NET_MANAGER', 1, 'Алматы',
 '{"street": "Абая 150", "building": "1"}',
 '{"monday": {"from": "09:00", "to": "21:00"}, "tuesday": {"from": "09:00", "to": "21:00"}}',
 true, 'test_seed'),
('POINT_NM_02', 'NetManager Точка 2', 'NET_MANAGER', 1, 'Алматы',
 '{"street": "Достык 200", "building": "2"}',
 '{"monday": {"from": "09:00", "to": "21:00"}, "tuesday": {"from": "09:00", "to": "21:00"}}',
 true, 'test_seed')
ON CONFLICT (code, network_code) DO NOTHING;

-- Создаем точки для сети SelfOwner
INSERT INTO points (code, name, network_code, category_id, city, address, schedule, active, updated_by) VALUES
('POINT_SO_01', 'SelfOwner Точка 1', 'NET_SELF_OWNER', 1, 'Астана',
 '{"street": "Кунаева 10", "building": "3"}',
 '{"monday": {"from": "09:00", "to": "21:00"}, "tuesday": {"from": "09:00", "to": "21:00"}}',
 true, 'test_seed'),
('POINT_SO_02', 'SelfOwner Точка 2', 'NET_SELF_OWNER', 1, 'Астана',
 '{"street": "Сарыарка 20", "building": "4"}',
 '{"monday": {"from": "09:00", "to": "21:00"}, "tuesday": {"from": "09:00", "to": "21:00"}}',
 true, 'test_seed')
ON CONFLICT (code, network_code) DO NOTHING;

-- Создаем точки для обычной сети (Manager)
INSERT INTO points (code, name, network_code, category_id, city, address, schedule, active, updated_by) VALUES
('POINT_MGR_01', 'Manager Точка 1', 'NET_REGULAR', 1, 'Шымкент',
 '{"street": "Тауке хана 5", "building": "5"}',
 '{"monday": {"from": "09:00", "to": "21:00"}, "tuesday": {"from": "09:00", "to": "21:00"}}',
 true, 'test_seed'),
('POINT_MGR_02', 'Manager Точка 2', 'NET_REGULAR', 1, 'Шымкент',
 '{"street": "Байтурсынова 15", "building": "6"}',
 '{"monday": {"from": "09:00", "to": "21:00"}, "tuesday": {"from": "09:00", "to": "21:00"}}',
 true, 'test_seed')
ON CONFLICT (code, network_code) DO NOTHING;

-- ============================================
-- МЕНЕДЖЕРЫ
-- ============================================

-- NetManager (может видеть всех пользователей своей сети)
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77001111111', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'net_manager', 'Алексей', 'Сетевой', 'NET_MANAGER', NULL, true, 'test_seed')
ON CONFLICT (phone) DO NOTHING;

-- SelfOwner (аналогично NetManager - может видеть всех пользователей своей сети)
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77002222222', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'self_owner', 'Владимир', 'Владелец', 'NET_SELF_OWNER', NULL, true, 'test_seed')
ON CONFLICT (phone) DO NOTHING;

-- Manager (может видеть только пользователей своей точки)
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77003333333', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'manager', 'Ирина', 'Менеджерова', 'NET_REGULAR', 'POINT_MGR_01', true, 'test_seed')
ON CONFLICT (phone) DO NOTHING;

-- ============================================
-- STAFF ПОЛЬЗОВАТЕЛИ ДЛЯ NETMANAGER
-- ============================================

-- 2 стафа из точки POINT_NM_01
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77001111101', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Айгуль', 'Мастерова', 'NET_MANAGER', 'POINT_NM_01', true, 'test_seed'),
('77001111102', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Асель', 'Стилистова', 'NET_MANAGER', 'POINT_NM_01', true, 'test_seed')
ON CONFLICT (phone) DO NOTHING;

-- 3 стафа из точки POINT_NM_02
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77001111201', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Дмитрий', 'Барбер', 'NET_MANAGER', 'POINT_NM_02', true, 'test_seed'),
('77001111202', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Евгений', 'Стрижкин', 'NET_MANAGER', 'POINT_NM_02', true, 'test_seed'),
('77001111203', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Максим', 'Прическов', 'NET_MANAGER', 'POINT_NM_02', true, 'test_seed')
ON CONFLICT (phone) DO NOTHING;

-- ============================================
-- STAFF ПОЛЬЗОВАТЕЛИ ДЛЯ SELFOWNER
-- ============================================

-- 2 стафа из точки POINT_SO_01
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77002222101', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Гульнара', 'Мастерица', 'NET_SELF_OWNER', 'POINT_SO_01', true, 'test_seed'),
('77002222102', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Жанна', 'Косметологова', 'NET_SELF_OWNER', 'POINT_SO_01', true, 'test_seed')
ON CONFLICT (phone) DO NOTHING;

-- 3 стафа из точки POINT_SO_02
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77002222201', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Инна', 'Стилистка', 'NET_SELF_OWNER', 'POINT_SO_02', true, 'test_seed'),
('77002222202', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Карина', 'Мастерская', 'NET_SELF_OWNER', 'POINT_SO_02', true, 'test_seed'),
('77002222203', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Лариса', 'Специалистова', 'NET_SELF_OWNER', 'POINT_SO_02', true, 'test_seed')
ON CONFLICT (phone) DO NOTHING;

-- ============================================
-- STAFF ПОЛЬЗОВАТЕЛИ ДЛЯ MANAGER
-- ============================================

-- 2 стафа из точки POINT_MGR_01
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77003333101', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Мария', 'Парикмахерова', 'NET_REGULAR', 'POINT_MGR_01', true, 'test_seed'),
('77003333102', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Наталья', 'Визажистова', 'NET_REGULAR', 'POINT_MGR_01', true, 'test_seed')
ON CONFLICT (phone) DO NOTHING;

-- 3 стафа из точки POINT_MGR_02
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77003333201', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Ольга', 'Маникюрша', 'NET_REGULAR', 'POINT_MGR_02', true, 'test_seed'),
('77003333202', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Павел', 'Стилист', 'NET_REGULAR', 'POINT_MGR_02', true, 'test_seed'),
('77003333203', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Роман', 'Колорист', 'NET_REGULAR', 'POINT_MGR_02', true, 'test_seed')
ON CONFLICT (phone) DO NOTHING;

-- ============================================
-- ДОПОЛНИТЕЛЬНЫЕ ТЕСТОВЫЕ ПОЛЬЗОВАТЕЛИ
-- ============================================

-- Обычный пользователь (клиент)
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77009999999', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'user', 'Тестовый', 'Клиент', NULL, NULL, true, 'test_seed')
ON CONFLICT (phone) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Удаляем тестовых пользователей
DELETE FROM users WHERE updated_by = 'test_seed';

-- Удаляем тестовые точки
DELETE FROM points WHERE updated_by = 'test_seed';

-- Удаляем тестовые сети
DELETE FROM networks WHERE updated_by = 'test_seed';

-- Удаляем тестовую категорию (опционально, т.к. может использоваться)
-- DELETE FROM point_categories WHERE updated_by = 'test_seed';

-- +goose StatementEnd