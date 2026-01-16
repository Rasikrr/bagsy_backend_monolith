-- +goose Up
-- +goose StatementBegin

-- ============================================
-- РЕАЛИСТИЧНЫЕ ТЕСТОВЫЕ ДАННЫЕ
-- ============================================
--
-- Пароль для всех пользователей: password123
-- Хеш: $2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe
--
-- ============================================

-- Категории точек (updated_by nullable)
INSERT INTO point_categories (name, description) VALUES
('Барбершоп', 'Мужская парикмахерская'),
('Салон красоты', 'Женский салон красоты'),
('Универсальный салон', 'Салон для мужчин и женщин');

-- ============================================
-- КАТЕГОРИИ УСЛУГ
-- ============================================

-- Категории услуг для барбершопов
INSERT INTO service_categories (name, description) VALUES
('Стрижки', 'Все виды мужских стрижек'),
('Бритье', 'Классическое и королевское бритье'),
('Уход за бородой', 'Стрижка, моделирование и окрашивание бороды'),
('Укладка и стайлинг', 'Укладка волос, стайлинг'),
('Окрашивание волос', 'Окрашивание и камуфлирование седины'),
('Маникюр', 'Все виды маникюра'),
('Педикюр', 'Все виды педикюра'),
('Косметология', 'Косметологические процедуры'),
('Массаж', 'Различные виды массажа');

-- Подкатегории для стрижек
INSERT INTO service_subcategories (service_category_id, name, description) VALUES
(1, 'Мужская классическая стрижка', 'Традиционная мужская стрижка'),
(1, 'Модельная стрижка', 'Современная модельная стрижка'),
(1, 'Детская стрижка', 'Стрижка для детей'),
(1, 'Женская стрижка', 'Женская стрижка любой сложности'),
(1, 'Стрижка машинкой', 'Короткая стрижка машинкой'),
-- Подкатегории для бритья
(2, 'Классическое бритье', 'Традиционное бритье опасной бритвой'),
(2, 'Королевское бритье', 'Бритье с горячими полотенцами и массажем'),
(2, 'Контурное бритье', 'Оформление контуров бороды'),
-- Подкатегории для ухода за бородой
(3, 'Стрижка бороды', 'Стрижка и оформление бороды'),
(3, 'Моделирование бороды', 'Создание формы бороды'),
(3, 'Окрашивание бороды', 'Окрашивание и тонирование бороды'),
(3, 'Коррекция бороды', 'Коррекция формы и длины'),
-- Подкатегории для укладки
(4, 'Укладка волос', 'Простая укладка'),
(4, 'Укладка феном', 'Профессиональная укладка феном'),
(4, 'Вечерняя укладка', 'Укладка для особых случаев'),
(4, 'Свадебная укладка', 'Свадебная прическа'),
-- Подкатегории для окрашивания волос
(5, 'Полное окрашивание', 'Окрашивание всех волос'),
(5, 'Мелирование', 'Окрашивание отдельных прядей'),
(5, 'Балаяж', 'Техника плавного перехода цвета'),
(5, 'Омбре', 'Градиентное окрашивание'),
(5, 'Камуфлирование седины', 'Маскировка седых волос'),
(5, 'Тонирование', 'Оттеночное окрашивание'),
-- Подкатегории для маникюра
(6, 'Классический маникюр', 'Традиционный маникюр'),
(6, 'Аппаратный маникюр', 'Маникюр с использованием аппарата'),
(6, 'Маникюр с покрытием', 'Маникюр с гель-лаком'),
(6, 'Наращивание ногтей', 'Наращивание гелем или акрилом'),
(6, 'Дизайн ногтей', 'Художественный дизайн'),
-- Подкатегории для педикюра
(7, 'Классический педикюр', 'Традиционный педикюр'),
(7, 'Аппаратный педикюр', 'Педикюр с использованием аппарата'),
(7, 'SPA-педикюр', 'Педикюр с уходовыми процедурами'),
(7, 'Медицинский педикюр', 'Лечебный педикюр'),
-- Подкатегории для косметологии
(8, 'Чистка лица', 'Механическая или ультразвуковая чистка'),
(8, 'Пилинг', 'Химический или механический пилинг'),
(8, 'Уходовая процедура', 'Увлажнение и питание кожи'),
(8, 'Массаж лица', 'Массаж для улучшения тонуса'),
(8, 'Маски для лица', 'Различные виды масок'),
-- Подкатегории для массажа
(9, 'Массаж головы', 'Расслабляющий массаж кожи головы'),
(9, 'Массаж шеи и воротниковой зоны', 'Снятие напряжения'),
(9, 'Классический массаж', 'Общий оздоровительный массаж');

-- ============================================
-- СЕТЬ 1: "Barbos" - Премиум барбершопы
-- ============================================

INSERT INTO networks (code, name, description, updated_by) VALUES
('BARBOS_KZ', 'Barbos', 'Сеть премиум барбершопов в Казахстане', 'realistic_seed');

-- Точки сети Barbos
INSERT INTO points (code, name, network_code, category_id, city, address, schedule, active, updated_by) VALUES
('BARBOS_ALMATY_ESENTAI', 'Barbos Esentai Mall', 'BARBOS_KZ', 1, 'Алматы',
 '{"coordinates": {"latitude": 43.2226, "longitude": 76.8512}, "street": "проспект Аль-Фараби, 77/8, 1 этаж", "city": "Алматы"}',
 '[{"week_day": 1, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 2, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 3, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 4, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 5, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T23:00:00Z", "all_day": false, "comment": ""}, {"week_day": 6, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T23:00:00Z", "all_day": false, "comment": ""}, {"week_day": 0, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}]',
 true, 'realistic_seed'),
('BARBOS_ALMATY_MEGA', 'Barbos Mega Center Almaty', 'BARBOS_KZ', 1, 'Алматы',
 '{"coordinates": {"latitude": 43.2073, "longitude": 76.6647}, "street": "Розыбакиева, 247А, 2 этаж", "city": "Алматы"}',
 '[{"week_day": 1, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 2, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 3, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 4, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 5, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T23:00:00Z", "all_day": false, "comment": ""}, {"week_day": 6, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T23:00:00Z", "all_day": false, "comment": ""}, {"week_day": 0, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}]',
 true, 'realistic_seed'),
('BARBOS_ASTANA_MEGA', 'Barbos Mega Silk Way', 'BARBOS_KZ', 1, 'Астана',
 '{"coordinates": {"latitude": 51.1282, "longitude": 71.4099}, "street": "Кабанбай батыра, 62, 1 этаж", "city": "Астана"}',
 '[{"week_day": 1, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 2, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 3, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 4, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 5, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T23:00:00Z", "all_day": false, "comment": ""}, {"week_day": 6, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T23:00:00Z", "all_day": false, "comment": ""}, {"week_day": 0, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}]',
 true, 'realistic_seed')
ON CONFLICT (code, network_code) DO NOTHING;

-- Управляющий сети Barbos
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77012345001', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'net_manager', 'Арман', 'Нурсултанов', 'BARBOS_KZ', NULL, true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- Менеджеры точек Barbos
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77012345101', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'manager', 'Ерлан', 'Қасымов', 'BARBOS_KZ', 'BARBOS_ALMATY_ESENTAI', true, 'realistic_seed'),
('77012345201', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'manager', 'Дмитрий', 'Павлов', 'BARBOS_KZ', 'BARBOS_ALMATY_MEGA', true, 'realistic_seed'),
('77012345301', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'manager', 'Нурлан', 'Сейтжанов', 'BARBOS_KZ', 'BARBOS_ASTANA_MEGA', true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- Барберы Barbos Esentai
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77012345111', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Данияр', 'Абдрахманов', 'BARBOS_KZ', 'BARBOS_ALMATY_ESENTAI', true, 'realistic_seed'),
('77012345112', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Максим', 'Ковалев', 'BARBOS_KZ', 'BARBOS_ALMATY_ESENTAI', true, 'realistic_seed'),
('77012345113', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Тимур', 'Жолдасов', 'BARBOS_KZ', 'BARBOS_ALMATY_ESENTAI', true, 'realistic_seed'),
('77012345114', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Александр', 'Сидоров', 'BARBOS_KZ', 'BARBOS_ALMATY_ESENTAI', true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- Барберы Barbos Mega Almaty
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77012345211', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Ержан', 'Омаров', 'BARBOS_KZ', 'BARBOS_ALMATY_MEGA', true, 'realistic_seed'),
('77012345212', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Руслан', 'Петров', 'BARBOS_KZ', 'BARBOS_ALMATY_MEGA', true, 'realistic_seed'),
('77012345213', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Серик', 'Ахметов', 'BARBOS_KZ', 'BARBOS_ALMATY_MEGA', true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- Барберы Barbos Astana
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77012345311', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Асхат', 'Досмагамбетов', 'BARBOS_KZ', 'BARBOS_ASTANA_MEGA', true, 'realistic_seed'),
('77012345312', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Игорь', 'Васильев', 'BARBOS_KZ', 'BARBOS_ASTANA_MEGA', true, 'realistic_seed'),
('77012345313', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Ринат', 'Сулейменов', 'BARBOS_KZ', 'BARBOS_ASTANA_MEGA', true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- ============================================
-- СЕТЬ 2: "Chic Style" - Женские салоны красоты
-- ============================================

INSERT INTO networks (code, name, description, updated_by) VALUES
('CHIC_STYLE_KZ', 'Chic Style', 'Сеть женских салонов красоты премиум-класса', 'realistic_seed')
ON CONFLICT (code) DO NOTHING;

-- Точки сети Chic Style
INSERT INTO points (code, name, network_code, category_id, city, address, schedule, active, updated_by) VALUES
('CHIC_ALMATY_FURMANOV', 'Chic Style Фурманова', 'CHIC_STYLE_KZ', 2, 'Алматы',
 '{"coordinates": {"latitude": 43.2567, "longitude": 76.9286}, "street": "улица Фурманова, 240, 1 этаж", "city": "Алматы"}',
 '[{"week_day": 1, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 2, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 3, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 4, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 5, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 6, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 0, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T20:00:00Z", "all_day": false, "comment": ""}]',
 true, 'realistic_seed'),
('CHIC_ALMATY_DOSTYK', 'Chic Style Достык', 'CHIC_STYLE_KZ', 2, 'Алматы',
 '{"coordinates": {"latitude": 43.2372, "longitude": 76.9453}, "street": "проспект Достык, 132, 2 этаж", "city": "Алматы"}',
 '[{"week_day": 1, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 2, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 3, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 4, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 5, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 6, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 0, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T20:00:00Z", "all_day": false, "comment": ""}]',
 true, 'realistic_seed')
ON CONFLICT (code, network_code) DO NOTHING;

-- Владелец сети Chic Style
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77013456001', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'self_owner', 'Айгерим', 'Сатпаева', 'CHIC_STYLE_KZ', NULL, true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- Менеджеры точек Chic Style
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77013456101', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'manager', 'Айнур', 'Абильдинова', 'CHIC_STYLE_KZ', 'CHIC_ALMATY_FURMANOV', true, 'realistic_seed'),
('77013456201', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'manager', 'Мадина', 'Есенова', 'CHIC_STYLE_KZ', 'CHIC_ALMATY_DOSTYK', true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- Мастера Chic Style Furmanov
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77013456111', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Асем', 'Нурланова', 'CHIC_STYLE_KZ', 'CHIC_ALMATY_FURMANOV', true, 'realistic_seed'),
('77013456112', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Дина', 'Жумабаева', 'CHIC_STYLE_KZ', 'CHIC_ALMATY_FURMANOV', true, 'realistic_seed'),
('77013456113', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Елена', 'Смирнова', 'CHIC_STYLE_KZ', 'CHIC_ALMATY_FURMANOV', true, 'realistic_seed'),
('77013456114', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Жанна', 'Бекмуратова', 'CHIC_STYLE_KZ', 'CHIC_ALMATY_FURMANOV', true, 'realistic_seed'),
('77013456115', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Индира', 'Касымова', 'CHIC_STYLE_KZ', 'CHIC_ALMATY_FURMANOV', true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- Мастера Chic Style Dostyk
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77013456211', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Камила', 'Оспанова', 'CHIC_STYLE_KZ', 'CHIC_ALMATY_DOSTYK', true, 'realistic_seed'),
('77013456212', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Лаура', 'Мухамбетова', 'CHIC_STYLE_KZ', 'CHIC_ALMATY_DOSTYK', true, 'realistic_seed'),
('77013456213', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Марина', 'Иванова', 'CHIC_STYLE_KZ', 'CHIC_ALMATY_DOSTYK', true, 'realistic_seed'),
('77013456214', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Назым', 'Қаржаубаева', 'CHIC_STYLE_KZ', 'CHIC_ALMATY_DOSTYK', true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- ============================================
-- СЕТЬ 3: "BroCode" - Сеть барбершопов среднего сегмента
-- ============================================

INSERT INTO networks (code, name, description, updated_by) VALUES
('BROCODE_KZ', 'BroCode', 'Городские барбершопы для современных мужчин', 'realistic_seed')
ON CONFLICT (code) DO NOTHING;

-- Точки сети BroCode
INSERT INTO points (code, name, network_code, category_id, city, address, schedule, active, updated_by) VALUES
('BROCODE_ALMATY_ABAY', 'BroCode на Абая', 'BROCODE_KZ', 1, 'Алматы',
 '{"coordinates": {"latitude": 43.2418, "longitude": 76.9011}, "street": "проспект Абая, 68, 1 этаж", "city": "Алматы"}',
 '[{"week_day": 1, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 2, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 3, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 4, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 5, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 6, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 0, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T20:00:00Z", "all_day": false, "comment": ""}]',
 true, 'realistic_seed'),
('BROCODE_ASTANA_RESPUBLIKA', 'BroCode Республика', 'BROCODE_KZ', 1, 'Астана',
 '{"coordinates": {"latitude": 51.1693, "longitude": 71.4491}, "street": "проспект Республики, 15, 1 этаж", "city": "Астана"}',
 '[{"week_day": 1, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 2, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 3, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 4, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 5, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 6, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 0, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T20:00:00Z", "all_day": false, "comment": ""}]',
 true, 'realistic_seed'),
('BROCODE_SHYMKENT_TAUKE', 'BroCode Тауке хана', 'BROCODE_KZ', 1, 'Шымкент',
 '{"coordinates": {"latitude": 42.3417, "longitude": 69.5901}, "street": "проспект Тауке хана, 32, 1 этаж", "city": "Шымкент"}',
 '[{"week_day": 1, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 2, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 3, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 4, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 5, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 6, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T22:00:00Z", "all_day": false, "comment": ""}, {"week_day": 0, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T20:00:00Z", "all_day": false, "comment": ""}]',
 true, 'realistic_seed')
ON CONFLICT (code, network_code) DO NOTHING;

-- Управляющий сети BroCode
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77014567001', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'net_manager', 'Бауыржан', 'Токаев', 'BROCODE_KZ', NULL, true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- Менеджеры точек BroCode
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77014567101', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'manager', 'Алмас', 'Сагындыков', 'BROCODE_KZ', 'BROCODE_ALMATY_ABAY', true, 'realistic_seed'),
('77014567201', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'manager', 'Олег', 'Кузнецов', 'BROCODE_KZ', 'BROCODE_ASTANA_RESPUBLIKA', true, 'realistic_seed'),
('77014567301', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'manager', 'Санжар', 'Мусаев', 'BROCODE_KZ', 'BROCODE_SHYMKENT_TAUKE', true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- Барберы BroCode Almaty
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77014567111', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Азамат', 'Досанов', 'BROCODE_KZ', 'BROCODE_ALMATY_ABAY', true, 'realistic_seed'),
('77014567112', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Батыр', 'Есимов', 'BROCODE_KZ', 'BROCODE_ALMATY_ABAY', true, 'realistic_seed'),
('77014567113', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Владислав', 'Морозов', 'BROCODE_KZ', 'BROCODE_ALMATY_ABAY', true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- Барберы BroCode Astana
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77014567211', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Дархан', 'Амангельдин', 'BROCODE_KZ', 'BROCODE_ASTANA_RESPUBLIKA', true, 'realistic_seed'),
('77014567212', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Евгений', 'Соколов', 'BROCODE_KZ', 'BROCODE_ASTANA_RESPUBLIKA', true, 'realistic_seed'),
('77014567213', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Жасулан', 'Қалиев', 'BROCODE_KZ', 'BROCODE_ASTANA_RESPUBLIKA', true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- Барберы BroCode Shymkent
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77014567311', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Кайрат', 'Төлеген', 'BROCODE_KZ', 'BROCODE_SHYMKENT_TAUKE', true, 'realistic_seed'),
('77014567312', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Марат', 'Ибрагимов', 'BROCODE_KZ', 'BROCODE_SHYMKENT_TAUKE', true, 'realistic_seed'),
('77014567313', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Нуржан', 'Сапаров', 'BROCODE_KZ', 'BROCODE_SHYMKENT_TAUKE', true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- ============================================
-- СЕТЬ 4: "Красотка" - Небольшая сеть салонов
-- ============================================

INSERT INTO networks (code, name, description, updated_by) VALUES
('KRASOTKA_KZ', 'Красотка', 'Доступные салоны красоты для всей семьи', 'realistic_seed')
ON CONFLICT (code) DO NOTHING;

-- Точки сети Красотка
INSERT INTO points (code, name, network_code, category_id, city, address, schedule, active, updated_by) VALUES
('KRASOTKA_ALMATY_ZHIBEK', 'Красотка на Жибек Жолы', 'KRASOTKA_KZ', 3, 'Алматы',
 '{"coordinates": {"latitude": 43.2630, "longitude": 76.9450}, "street": "проспект Жибек Жолы, 53, 1 этаж", "city": "Алматы"}',
 '[{"week_day": 1, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T20:00:00Z", "all_day": false, "comment": ""}, {"week_day": 2, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T20:00:00Z", "all_day": false, "comment": ""}, {"week_day": 3, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T20:00:00Z", "all_day": false, "comment": ""}, {"week_day": 4, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T20:00:00Z", "all_day": false, "comment": ""}, {"week_day": 5, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 6, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 0, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T19:00:00Z", "all_day": false, "comment": ""}]',
 true, 'realistic_seed'),
('KRASOTKA_ALMATY_ALATAU', 'Красотка Алатауский', 'KRASOTKA_KZ', 3, 'Алматы',
 '{"coordinates": {"latitude": 43.1894, "longitude": 76.9132}, "street": "микрорайон Алатау, 7А, 1 этаж", "city": "Алматы"}',
 '[{"week_day": 1, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T20:00:00Z", "all_day": false, "comment": ""}, {"week_day": 2, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T20:00:00Z", "all_day": false, "comment": ""}, {"week_day": 3, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T20:00:00Z", "all_day": false, "comment": ""}, {"week_day": 4, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T20:00:00Z", "all_day": false, "comment": ""}, {"week_day": 5, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 6, "open": "2000-01-01T09:00:00Z", "close": "2000-01-01T21:00:00Z", "all_day": false, "comment": ""}, {"week_day": 0, "open": "2000-01-01T10:00:00Z", "close": "2000-01-01T19:00:00Z", "all_day": false, "comment": ""}]',
 true, 'realistic_seed')
ON CONFLICT (code, network_code) DO NOTHING;

-- Владелец сети Красотка
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77015678001', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'self_owner', 'Гульнара', 'Абишева', 'KRASOTKA_KZ', NULL, true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- Менеджеры точек Красотка
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77015678101', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'manager', 'Сауле', 'Бейсембаева', 'KRASOTKA_KZ', 'KRASOTKA_ALMATY_ZHIBEK', true, 'realistic_seed'),
('77015678201', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'manager', 'Татьяна', 'Новикова', 'KRASOTKA_KZ', 'KRASOTKA_ALMATY_ALATAU', true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- Мастера Красотка Жибек Жолы
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77015678111', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Алия', 'Ержанова', 'KRASOTKA_KZ', 'KRASOTKA_ALMATY_ZHIBEK', true, 'realistic_seed'),
('77015678112', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Балжан', 'Қадырова', 'KRASOTKA_KZ', 'KRASOTKA_ALMATY_ZHIBEK', true, 'realistic_seed'),
('77015678113', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Виктория', 'Сергеева', 'KRASOTKA_KZ', 'KRASOTKA_ALMATY_ZHIBEK', true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- Мастера Красотка Алатау
INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77015678211', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Динара', 'Смагулова', 'KRASOTKA_KZ', 'KRASOTKA_ALMATY_ALATAU', true, 'realistic_seed'),
('77015678212', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Екатерина', 'Волкова', 'KRASOTKA_KZ', 'KRASOTKA_ALMATY_ALATAU', true, 'realistic_seed'),
('77015678213', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'staff', 'Жанар', 'Туғанбаева', 'KRASOTKA_KZ', 'KRASOTKA_ALMATY_ALATAU', true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- ============================================
-- ТЕСТОВЫЕ КЛИЕНТЫ
-- ============================================

INSERT INTO users (phone, password, role, name, surname, network_code, point_code, active, updated_by) VALUES
('77016789001', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'user', 'Айдос', 'Мамедов', NULL, NULL, true, 'realistic_seed'),
('77016789002', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'user', 'Мария', 'Ким', NULL, NULL, true, 'realistic_seed'),
('77016789003', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'user', 'Ерболат', 'Сейтов', NULL, NULL, true, 'realistic_seed'),
('77016789004', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'user', 'Анна', 'Петрова', NULL, NULL, true, 'realistic_seed'),
('77016789005', '$2a$10$Se9s4eV/CB8e12VYuPmT/.Zyw4cBim8FEeItq1aC60PXGJguyNfNe', 'user', 'Нурбол', 'Алдабергенов', NULL, NULL, true, 'realistic_seed')
ON CONFLICT (phone) DO NOTHING;

-- ============================================
-- УСЛУГИ ДЛЯ ТОЧЕК
-- ============================================

-- Услуги для барбершопов Barbos
-- BARBOS_ALMATY_ESENTAI
INSERT INTO services (point_code, category_id, subcategory_id, name, description, duration_minutes, active) VALUES
('BARBOS_ALMATY_ESENTAI', 1, 1, 'Мужская классическая стрижка', 'Классическая мужская стрижка ножницами и машинкой', 40, true),
('BARBOS_ALMATY_ESENTAI', 1, 2, 'Модельная стрижка', 'Современная модельная стрижка по вашему желанию', 50, true),
('BARBOS_ALMATY_ESENTAI', 1, 3, 'Детская стрижка', 'Стрижка для мальчиков до 12 лет', 30, true),
('BARBOS_ALMATY_ESENTAI', 2, 6, 'Классическое бритье', 'Традиционное бритье опасной бритвой с горячими полотенцами', 40, true),
('BARBOS_ALMATY_ESENTAI', 2, 7, 'Королевское бритье', 'Бритье премиум-класса с массажем лица и ароматерапией', 60, true),
('BARBOS_ALMATY_ESENTAI', 3, 9, 'Стрижка бороды', 'Оформление и стрижка бороды', 30, true),
('BARBOS_ALMATY_ESENTAI', 3, 10, 'Моделирование бороды', 'Создание формы бороды по вашему желанию', 40, true),
('BARBOS_ALMATY_ESENTAI', 3, 11, 'Окрашивание бороды', 'Окрашивание или камуфлирование седины', 45, true),
('BARBOS_ALMATY_ESENTAI', 4, 13, 'Укладка волос', 'Укладка с использованием профессиональных средств', 20, true),
('BARBOS_ALMATY_ESENTAI', 5, 17, 'Камуфлирование седины', 'Маскировка седых волос натуральными тонами', 50, true);

-- BARBOS_ALMATY_MEGA
INSERT INTO services (point_code, category_id, subcategory_id, name, description, duration_minutes, active) VALUES
('BARBOS_ALMATY_MEGA', 1, 1, 'Мужская классическая стрижка', 'Классическая мужская стрижка', 40, true),
('BARBOS_ALMATY_MEGA', 1, 2, 'Модельная стрижка', 'Современная стрижка любой сложности', 50, true),
('BARBOS_ALMATY_MEGA', 1, 3, 'Детская стрижка', 'Стрижка для детей', 30, true),
('BARBOS_ALMATY_MEGA', 2, 6, 'Классическое бритье', 'Бритье опасной бритвой', 40, true),
('BARBOS_ALMATY_MEGA', 3, 9, 'Стрижка бороды', 'Стрижка и оформление бороды', 30, true),
('BARBOS_ALMATY_MEGA', 3, 10, 'Моделирование бороды', 'Создание и коррекция формы', 40, true),
('BARBOS_ALMATY_MEGA', 4, 13, 'Укладка волос', 'Профессиональная укладка', 20, true),
('BARBOS_ALMATY_MEGA', 5, 17, 'Камуфлирование седины', 'Маскировка седины', 50, true);

-- BARBOS_ASTANA_MEGA
INSERT INTO services (point_code, category_id, subcategory_id, name, description, duration_minutes, active) VALUES
('BARBOS_ASTANA_MEGA', 1, 1, 'Мужская классическая стрижка', 'Классическая стрижка', 40, true),
('BARBOS_ASTANA_MEGA', 1, 2, 'Модельная стрижка', 'Модельная стрижка', 50, true),
('BARBOS_ASTANA_MEGA', 1, 3, 'Детская стрижка', 'Детская стрижка', 30, true),
('BARBOS_ASTANA_MEGA', 2, 6, 'Классическое бритье', 'Классическое бритье', 40, true),
('BARBOS_ASTANA_MEGA', 2, 7, 'Королевское бритье', 'Премиум бритье', 60, true),
('BARBOS_ASTANA_MEGA', 3, 9, 'Стрижка бороды', 'Стрижка бороды', 30, true),
('BARBOS_ASTANA_MEGA', 3, 11, 'Окрашивание бороды', 'Окрашивание', 45, true);

-- Услуги для салонов Chic Style
-- CHIC_ALMATY_FURMANOV
INSERT INTO services (point_code, category_id, subcategory_id, name, description, duration_minutes, active) VALUES
('CHIC_ALMATY_FURMANOV', 1, 4, 'Женская стрижка', 'Стрижка любой сложности', 60, true),
('CHIC_ALMATY_FURMANOV', 4, 14, 'Укладка феном', 'Профессиональная укладка', 40, true),
('CHIC_ALMATY_FURMANOV', 4, 15, 'Вечерняя укладка', 'Укладка для особых случаев', 90, true),
('CHIC_ALMATY_FURMANOV', 4, 16, 'Свадебная укладка', 'Свадебная прическа', 120, true),
('CHIC_ALMATY_FURMANOV', 5, 19, 'Полное окрашивание', 'Окрашивание всех волос', 120, true),
('CHIC_ALMATY_FURMANOV', 5, 20, 'Мелирование', 'Мелирование прядей', 150, true),
('CHIC_ALMATY_FURMANOV', 5, 21, 'Балаяж', 'Техника балаяж', 180, true),
('CHIC_ALMATY_FURMANOV', 5, 22, 'Омбре', 'Градиентное окрашивание', 180, true),
('CHIC_ALMATY_FURMANOV', 6, 25, 'Классический маникюр', 'Традиционный маникюр', 60, true),
('CHIC_ALMATY_FURMANOV', 6, 26, 'Аппаратный маникюр', 'Маникюр аппаратом', 60, true),
('CHIC_ALMATY_FURMANOV', 6, 27, 'Маникюр с покрытием', 'Маникюр с гель-лаком', 90, true),
('CHIC_ALMATY_FURMANOV', 7, 29, 'Классический педикюр', 'Традиционный педикюр', 60, true),
('CHIC_ALMATY_FURMANOV', 7, 30, 'Аппаратный педикюр', 'Педикюр аппаратом', 60, true),
('CHIC_ALMATY_FURMANOV', 7, 31, 'SPA-педикюр', 'Педикюр с уходом', 90, true),
('CHIC_ALMATY_FURMANOV', 8, 33, 'Чистка лица', 'Глубокая чистка', 60, true),
('CHIC_ALMATY_FURMANOV', 8, 35, 'Уходовая процедура', 'Увлажнение и питание', 60, true);

-- CHIC_ALMATY_DOSTYK
INSERT INTO services (point_code, category_id, subcategory_id, name, description, duration_minutes, active) VALUES
('CHIC_ALMATY_DOSTYK', 1, 4, 'Женская стрижка', 'Стрижка любой сложности', 60, true),
('CHIC_ALMATY_DOSTYK', 4, 14, 'Укладка феном', 'Укладка', 40, true),
('CHIC_ALMATY_DOSTYK', 4, 15, 'Вечерняя укладка', 'Вечерняя укладка', 90, true),
('CHIC_ALMATY_DOSTYK', 5, 19, 'Полное окрашивание', 'Окрашивание', 120, true),
('CHIC_ALMATY_DOSTYK', 5, 21, 'Балаяж', 'Балаяж', 180, true),
('CHIC_ALMATY_DOSTYK', 6, 26, 'Аппаратный маникюр', 'Маникюр', 60, true),
('CHIC_ALMATY_DOSTYK', 6, 27, 'Маникюр с покрытием', 'С гель-лаком', 90, true),
('CHIC_ALMATY_DOSTYK', 7, 30, 'Аппаратный педикюр', 'Педикюр', 60, true),
('CHIC_ALMATY_DOSTYK', 8, 33, 'Чистка лица', 'Чистка', 60, true);

-- Услуги для BroCode
-- BROCODE_ALMATY_ABAY
INSERT INTO services (point_code, category_id, subcategory_id, name, description, duration_minutes, active) VALUES
('BROCODE_ALMATY_ABAY', 1, 1, 'Мужская классическая стрижка', 'Классика', 40, true),
('BROCODE_ALMATY_ABAY', 1, 2, 'Модельная стрижка', 'Модельная', 50, true),
('BROCODE_ALMATY_ABAY', 1, 3, 'Детская стрижка', 'Детская', 30, true),
('BROCODE_ALMATY_ABAY', 2, 6, 'Классическое бритье', 'Бритье', 40, true),
('BROCODE_ALMATY_ABAY', 3, 9, 'Стрижка бороды', 'Борода', 30, true),
('BROCODE_ALMATY_ABAY', 4, 13, 'Укладка волос', 'Укладка', 20, true);

-- BROCODE_ASTANA_RESPUBLIKA
INSERT INTO services (point_code, category_id, subcategory_id, name, description, duration_minutes, active) VALUES
('BROCODE_ASTANA_RESPUBLIKA', 1, 1, 'Мужская классическая стрижка', 'Классика', 40, true),
('BROCODE_ASTANA_RESPUBLIKA', 1, 2, 'Модельная стрижка', 'Модельная', 50, true),
('BROCODE_ASTANA_RESPUBLIKA', 2, 6, 'Классическое бритье', 'Бритье', 40, true),
('BROCODE_ASTANA_RESPUBLIKA', 3, 9, 'Стрижка бороды', 'Борода', 30, true);

-- BROCODE_SHYMKENT_TAUKE
INSERT INTO services (point_code, category_id, subcategory_id, name, description, duration_minutes, active) VALUES
('BROCODE_SHYMKENT_TAUKE', 1, 1, 'Мужская классическая стрижка', 'Классика', 40, true),
('BROCODE_SHYMKENT_TAUKE', 1, 2, 'Модельная стрижка', 'Модельная', 50, true),
('BROCODE_SHYMKENT_TAUKE', 1, 3, 'Детская стрижка', 'Детская', 30, true),
('BROCODE_SHYMKENT_TAUKE', 2, 6, 'Классическое бритье', 'Бритье', 40, true),
('BROCODE_SHYMKENT_TAUKE', 3, 9, 'Стрижка бороды', 'Борода', 30, true);

-- Услуги для Krasotka (универсальные салоны)
-- KRASOTKA_ALMATY_ZHIBEK
INSERT INTO services (point_code, category_id, subcategory_id, name, description, duration_minutes, active) VALUES
('KRASOTKA_ALMATY_ZHIBEK', 1, 1, 'Мужская стрижка', 'Мужская', 40, true),
('KRASOTKA_ALMATY_ZHIBEK', 1, 4, 'Женская стрижка', 'Женская', 60, true),
('KRASOTKA_ALMATY_ZHIBEK', 1, 3, 'Детская стрижка', 'Детская', 30, true),
('KRASOTKA_ALMATY_ZHIBEK', 4, 13, 'Укладка волос', 'Укладка', 30, true),
('KRASOTKA_ALMATY_ZHIBEK', 5, 19, 'Полное окрашивание', 'Окрашивание', 120, true),
('KRASOTKA_ALMATY_ZHIBEK', 6, 25, 'Классический маникюр', 'Маникюр', 60, true),
('KRASOTKA_ALMATY_ZHIBEK', 7, 29, 'Классический педикюр', 'Педикюр', 60, true);

-- KRASOTKA_ALMATY_ALATAU
INSERT INTO services (point_code, category_id, subcategory_id, name, description, duration_minutes, active) VALUES
('KRASOTKA_ALMATY_ALATAU', 1, 1, 'Мужская стрижка', 'Мужская', 40, true),
('KRASOTKA_ALMATY_ALATAU', 1, 4, 'Женская стрижка', 'Женская', 60, true),
('KRASOTKA_ALMATY_ALATAU', 1, 3, 'Детская стрижка', 'Детская', 30, true),
('KRASOTKA_ALMATY_ALATAU', 4, 13, 'Укладка волос', 'Укладка', 30, true),
('KRASOTKA_ALMATY_ALATAU', 5, 19, 'Полное окрашивание', 'Окрашивание', 120, true),
('KRASOTKA_ALMATY_ALATAU', 6, 26, 'Аппаратный маникюр', 'Маникюр', 60, true),
('KRASOTKA_ALMATY_ALATAU', 7, 30, 'Аппаратный педикюр', 'Педикюр', 60, true);

-- ============================================
-- ПРИВЯЗКА МАСТЕРОВ К УСЛУГАМ (MASTER_SERVICES)
-- ============================================

-- BARBOS_ALMATY_ESENTAI - 4 барбера, все услуги точки
INSERT INTO master_services (master_phone, service_id, price, active, updated_by)
SELECT data.master_phone, s.id, data.price, true, 'realistic_seed'
FROM (VALUES
    -- Данияр Абдрахманов (77012345111)
    ('77012345111', 'BARBOS_ALMATY_ESENTAI', 'Мужская классическая стрижка', 8000.00),
    ('77012345111', 'BARBOS_ALMATY_ESENTAI', 'Модельная стрижка', 10000.00),
    ('77012345111', 'BARBOS_ALMATY_ESENTAI', 'Детская стрижка', 6000.00),
    ('77012345111', 'BARBOS_ALMATY_ESENTAI', 'Классическое бритье', 6000.00),
    ('77012345111', 'BARBOS_ALMATY_ESENTAI', 'Стрижка бороды', 5000.00),
    -- Максим Ковалев (77012345112)
    ('77012345112', 'BARBOS_ALMATY_ESENTAI', 'Мужская классическая стрижка', 9000.00),
    ('77012345112', 'BARBOS_ALMATY_ESENTAI', 'Модельная стрижка', 11000.00),
    ('77012345112', 'BARBOS_ALMATY_ESENTAI', 'Королевское бритье', 8000.00),
    ('77012345112', 'BARBOS_ALMATY_ESENTAI', 'Моделирование бороды', 6000.00),
    -- Тимур Жолдасов (77012345113)
    ('77012345113', 'BARBOS_ALMATY_ESENTAI', 'Мужская классическая стрижка', 8500.00),
    ('77012345113', 'BARBOS_ALMATY_ESENTAI', 'Модельная стрижка', 10500.00),
    ('77012345113', 'BARBOS_ALMATY_ESENTAI', 'Детская стрижка', 5500.00),
    ('77012345113', 'BARBOS_ALMATY_ESENTAI', 'Окрашивание бороды', 7000.00),
    -- Александр Сидоров (77012345114)
    ('77012345114', 'BARBOS_ALMATY_ESENTAI', 'Мужская классическая стрижка', 9500.00),
    ('77012345114', 'BARBOS_ALMATY_ESENTAI', 'Модельная стрижка', 12000.00),
    ('77012345114', 'BARBOS_ALMATY_ESENTAI', 'Королевское бритье', 8500.00),
    ('77012345114', 'BARBOS_ALMATY_ESENTAI', 'Стрижка бороды', 5500.00)
) AS data(master_phone, point_code, service_name, price)
JOIN services s ON s.point_code = data.point_code AND s.name = data.service_name;

-- BARBOS_ALMATY_MEGA - 3 барбера
INSERT INTO master_services (master_phone, service_id, price, active, updated_by)
SELECT data.master_phone, s.id, data.price, true, 'realistic_seed'
FROM (VALUES
    -- Ержан Омаров (77012345211)
    ('77012345211', 'BARBOS_ALMATY_MEGA', 'Мужская классическая стрижка', 7500.00),
    ('77012345211', 'BARBOS_ALMATY_MEGA', 'Модельная стрижка', 9500.00),
    ('77012345211', 'BARBOS_ALMATY_MEGA', 'Детская стрижка', 5000.00),
    ('77012345211', 'BARBOS_ALMATY_MEGA', 'Классическое бритье', 5500.00),
    -- Руслан Петров (77012345212)
    ('77012345212', 'BARBOS_ALMATY_MEGA', 'Мужская классическая стрижка', 8000.00),
    ('77012345212', 'BARBOS_ALMATY_MEGA', 'Модельная стрижка', 10000.00),
    ('77012345212', 'BARBOS_ALMATY_MEGA', 'Стрижка бороды', 4500.00),
    ('77012345212', 'BARBOS_ALMATY_MEGA', 'Королевское бритье', 7500.00),
    -- Серик Ахметов (77012345213)
    ('77012345213', 'BARBOS_ALMATY_MEGA', 'Мужская классическая стрижка', 7000.00),
    ('77012345213', 'BARBOS_ALMATY_MEGA', 'Модельная стрижка', 9000.00),
    ('77012345213', 'BARBOS_ALMATY_MEGA', 'Детская стрижка', 4500.00),
    ('77012345213', 'BARBOS_ALMATY_MEGA', 'Моделирование бороды', 5000.00)
) AS data(master_phone, point_code, service_name, price)
JOIN services s ON s.point_code = data.point_code AND s.name = data.service_name;

-- BARBOS_ASTANA_MEGA - 3 барбера
INSERT INTO master_services (master_phone, service_id, price, active, updated_by)
SELECT data.master_phone, s.id, data.price, true, 'realistic_seed'
FROM (VALUES
    -- Асхат Досмагамбетов (77012345311)
    ('77012345311', 'BARBOS_ASTANA_MEGA', 'Мужская классическая стрижка', 7500.00),
    ('77012345311', 'BARBOS_ASTANA_MEGA', 'Модельная стрижка', 9500.00),
    ('77012345311', 'BARBOS_ASTANA_MEGA', 'Классическое бритье', 5500.00),
    ('77012345311', 'BARBOS_ASTANA_MEGA', 'Массаж головы', 3500.00),
    -- Игорь Васильев (77012345312)
    ('77012345312', 'BARBOS_ASTANA_MEGA', 'Мужская классическая стрижка', 8000.00),
    ('77012345312', 'BARBOS_ASTANA_MEGA', 'Модельная стрижка', 10000.00),
    ('77012345312', 'BARBOS_ASTANA_MEGA', 'Королевское бритье', 7000.00),
    ('77012345312', 'BARBOS_ASTANA_MEGA', 'Стрижка бороды', 4500.00),
    -- Ринат Сулейменов (77012345313)
    ('77012345313', 'BARBOS_ASTANA_MEGA', 'Мужская классическая стрижка', 7000.00),
    ('77012345313', 'BARBOS_ASTANA_MEGA', 'Модельная стрижка', 9000.00),
    ('77012345313', 'BARBOS_ASTANA_MEGA', 'Детская стрижка', 5000.00),
    ('77012345313', 'BARBOS_ASTANA_MEGA', 'Моделирование бороды', 5500.00)
) AS data(master_phone, point_code, service_name, price)
JOIN services s ON s.point_code = data.point_code AND s.name = data.service_name;

-- CHIC_ALMATY_FURMANOV - 5 мастеров, женские услуги
INSERT INTO master_services (master_phone, service_id, price, active, updated_by)
SELECT data.master_phone, s.id, data.price, true, 'realistic_seed'
FROM (VALUES
    -- Асем Нурланова (77013456111) - стрижки и укладки
    ('77013456111', 'CHIC_ALMATY_FURMANOV', 'Женская стрижка', 12000.00),
    ('77013456111', 'CHIC_ALMATY_FURMANOV', 'Укладка феном', 7000.00),
    ('77013456111', 'CHIC_ALMATY_FURMANOV', 'Вечерняя укладка', 10000.00),
    ('77013456111', 'CHIC_ALMATY_FURMANOV', 'Свадебная прическа', 15000.00),
    -- Дина Жумабаева (77013456112) - окрашивание
    ('77013456112', 'CHIC_ALMATY_FURMANOV', 'Женская стрижка', 11000.00),
    ('77013456112', 'CHIC_ALMATY_FURMANOV', 'Полное окрашивание', 20000.00),
    ('77013456112', 'CHIC_ALMATY_FURMANOV', 'Мелирование', 22000.00),
    ('77013456112', 'CHIC_ALMATY_FURMANOV', 'Балаяж', 25000.00),
    ('77013456112', 'CHIC_ALMATY_FURMANOV', 'Омбре', 24000.00),
    -- Елена Смирнова (77013456113) - маникюр
    ('77013456113', 'CHIC_ALMATY_FURMANOV', 'Классический маникюр', 6000.00),
    ('77013456113', 'CHIC_ALMATY_FURMANOV', 'Аппаратный маникюр', 7000.00),
    ('77013456113', 'CHIC_ALMATY_FURMANOV', 'Маникюр с покрытием', 9000.00),
    -- Жанна Бекмуратова (77013456114) - педикюр
    ('77013456114', 'CHIC_ALMATY_FURMANOV', 'Классический педикюр', 8000.00),
    ('77013456114', 'CHIC_ALMATY_FURMANOV', 'Аппаратный педикюр', 9000.00),
    ('77013456114', 'CHIC_ALMATY_FURMANOV', 'SPA-педикюр', 12000.00),
    -- Индира Касымова (77013456115) - косметология
    ('77013456115', 'CHIC_ALMATY_FURMANOV', 'Чистка лица', 12000.00),
    ('77013456115', 'CHIC_ALMATY_FURMANOV', 'Уходовая процедура', 10000.00)
) AS data(master_phone, point_code, service_name, price)
JOIN services s ON s.point_code = data.point_code AND s.name = data.service_name;

-- CHIC_ALMATY_DOSTYK - 4 мастера
INSERT INTO master_services (master_phone, service_id, price, active, updated_by)
SELECT data.master_phone, s.id, data.price, true, 'realistic_seed'
FROM (VALUES
    -- Камила Оспанова (77013456211) - стрижки и укладки
    ('77013456211', 'CHIC_ALMATY_DOSTYK', 'Женская стрижка', 11000.00),
    ('77013456211', 'CHIC_ALMATY_DOSTYK', 'Укладка феном', 6500.00),
    ('77013456211', 'CHIC_ALMATY_DOSTYK', 'Вечерняя укладка', 9500.00),
    -- Лаура Мухамбетова (77013456212) - окрашивание
    ('77013456212', 'CHIC_ALMATY_DOSTYK', 'Полное окрашивание', 18000.00),
    ('77013456212', 'CHIC_ALMATY_DOSTYK', 'Балаяж', 23000.00),
    -- Марина Иванова (77013456213) - маникюр
    ('77013456213', 'CHIC_ALMATY_DOSTYK', 'Аппаратный маникюр', 6500.00),
    ('77013456213', 'CHIC_ALMATY_DOSTYK', 'Маникюр с покрытием', 8500.00),
    -- Назым Қаржаубаева (77013456214) - педикюр и косметология
    ('77013456214', 'CHIC_ALMATY_DOSTYK', 'Аппаратный педикюр', 8500.00),
    ('77013456214', 'CHIC_ALMATY_DOSTYK', 'Чистка лица', 11000.00)
) AS data(master_phone, point_code, service_name, price)
JOIN services s ON s.point_code = data.point_code AND s.name = data.service_name;

-- BROCODE_ALMATY_ABAY - 3 барбера
INSERT INTO master_services (master_phone, service_id, price, active, updated_by)
SELECT data.master_phone, s.id, data.price, true, 'realistic_seed'
FROM (VALUES
    -- Азамат Досанов (77014567111)
    ('77014567111', 'BROCODE_ALMATY_ABAY', 'Мужская классическая стрижка', 6500.00),
    ('77014567111', 'BROCODE_ALMATY_ABAY', 'Модельная стрижка', 8500.00),
    ('77014567111', 'BROCODE_ALMATY_ABAY', 'Классическое бритье', 5000.00),
    -- Батыр Есимов (77014567112)
    ('77014567112', 'BROCODE_ALMATY_ABAY', 'Мужская классическая стрижка', 7000.00),
    ('77014567112', 'BROCODE_ALMATY_ABAY', 'Модельная стрижка', 9000.00),
    ('77014567112', 'BROCODE_ALMATY_ABAY', 'Детская стрижка', 4500.00),
    ('77014567112', 'BROCODE_ALMATY_ABAY', 'Стрижка бороды', 4000.00),
    -- Владислав Морозов (77014567113)
    ('77014567113', 'BROCODE_ALMATY_ABAY', 'Мужская классическая стрижка', 6000.00),
    ('77014567113', 'BROCODE_ALMATY_ABAY', 'Модельная стрижка', 8000.00),
    ('77014567113', 'BROCODE_ALMATY_ABAY', 'Укладка волос', 3500.00)
) AS data(master_phone, point_code, service_name, price)
JOIN services s ON s.point_code = data.point_code AND s.name = data.service_name;

-- BROCODE_ASTANA_RESPUBLIKA - 3 барбера
INSERT INTO master_services (master_phone, service_id, price, active, updated_by)
SELECT data.master_phone, s.id, data.price, true, 'realistic_seed'
FROM (VALUES
    -- Дархан Амангельдин (77014567211)
    ('77014567211', 'BROCODE_ASTANA_RESPUBLIKA', 'Мужская классическая стрижка', 6000.00),
    ('77014567211', 'BROCODE_ASTANA_RESPUBLIKA', 'Модельная стрижка', 8000.00),
    -- Евгений Соколов (77014567212)
    ('77014567212', 'BROCODE_ASTANA_RESPUBLIKA', 'Мужская классическая стрижка', 6500.00),
    ('77014567212', 'BROCODE_ASTANA_RESPUBLIKA', 'Классическое бритье', 4500.00),
    ('77014567212', 'BROCODE_ASTANA_RESPUBLIKA', 'Стрижка бороды', 3500.00),
    -- Жасулан Қалиев (77014567213)
    ('77014567213', 'BROCODE_ASTANA_RESPUBLIKA', 'Модельная стрижка', 7500.00),
    ('77014567213', 'BROCODE_ASTANA_RESPUBLIKA', 'Стрижка бороды', 4000.00)
) AS data(master_phone, point_code, service_name, price)
JOIN services s ON s.point_code = data.point_code AND s.name = data.service_name;

-- BROCODE_SHYMKENT_TAUKE - 3 барбера
INSERT INTO master_services (master_phone, service_id, price, active, updated_by)
SELECT data.master_phone, s.id, data.price, true, 'realistic_seed'
FROM (VALUES
    -- Кайрат Төлеген (77014567311)
    ('77014567311', 'BROCODE_SHYMKENT_TAUKE', 'Мужская классическая стрижка', 5500.00),
    ('77014567311', 'BROCODE_SHYMKENT_TAUKE', 'Модельная стрижка', 7500.00),
    ('77014567311', 'BROCODE_SHYMKENT_TAUKE', 'Детская стрижка', 4000.00),
    -- Марат Ибрагимов (77014567312)
    ('77014567312', 'BROCODE_SHYMKENT_TAUKE', 'Мужская классическая стрижка', 6000.00),
    ('77014567312', 'BROCODE_SHYMKENT_TAUKE', 'Классическое бритье', 4000.00),
    ('77014567312', 'BROCODE_SHYMKENT_TAUKE', 'Стрижка бороды', 3500.00),
    -- Нуржан Сапаров (77014567313)
    ('77014567313', 'BROCODE_SHYMKENT_TAUKE', 'Мужская классическая стрижка', 5000.00),
    ('77014567313', 'BROCODE_SHYMKENT_TAUKE', 'Модельная стрижка', 7000.00),
    ('77014567313', 'BROCODE_SHYMKENT_TAUKE', 'Детская стрижка', 3500.00)
) AS data(master_phone, point_code, service_name, price)
JOIN services s ON s.point_code = data.point_code AND s.name = data.service_name;

-- KRASOTKA_ALMATY_ZHIBEK - 3 мастера (универсальный салон)
INSERT INTO master_services (master_phone, service_id, price, active, updated_by)
SELECT data.master_phone, s.id, data.price, true, 'realistic_seed'
FROM (VALUES
    -- Алия Ержанова (77015678111) - стрижки
    ('77015678111', 'KRASOTKA_ALMATY_ZHIBEK', 'Мужская стрижка', 6000.00),
    ('77015678111', 'KRASOTKA_ALMATY_ZHIBEK', 'Женская стрижка', 10000.00),
    ('77015678111', 'KRASOTKA_ALMATY_ZHIBEK', 'Детская стрижка', 4000.00),
    ('77015678111', 'KRASOTKA_ALMATY_ZHIBEK', 'Укладка волос', 5000.00),
    -- Балжан Қадырова (77015678112) - окрашивание и маникюр
    ('77015678112', 'KRASOTKA_ALMATY_ZHIBEK', 'Полное окрашивание', 15000.00),
    ('77015678112', 'KRASOTKA_ALMATY_ZHIBEK', 'Классический маникюр', 5500.00),
    -- Виктория Сергеева (77015678113) - педикюр
    ('77015678113', 'KRASOTKA_ALMATY_ZHIBEK', 'Классический педикюр', 7000.00)
) AS data(master_phone, point_code, service_name, price)
JOIN services s ON s.point_code = data.point_code AND s.name = data.service_name;

-- KRASOTKA_ALMATY_ALATAU - 3 мастера
INSERT INTO master_services (master_phone, service_id, price, active, updated_by)
SELECT data.master_phone, s.id, data.price, true, 'realistic_seed'
FROM (VALUES
    -- Динара Смагулова (77015678211) - стрижки
    ('77015678211', 'KRASOTKA_ALMATY_ALATAU', 'Мужская стрижка', 5500.00),
    ('77015678211', 'KRASOTKA_ALMATY_ALATAU', 'Женская стрижка', 9500.00),
    ('77015678211', 'KRASOTKA_ALMATY_ALATAU', 'Детская стрижка', 3500.00),
    ('77015678211', 'KRASOTKA_ALMATY_ALATAU', 'Укладка волос', 4500.00),
    -- Екатерина Волкова (77015678212) - окрашивание и маникюр
    ('77015678212', 'KRASOTKA_ALMATY_ALATAU', 'Полное окрашивание', 14000.00),
    ('77015678212', 'KRASOTKA_ALMATY_ALATAU', 'Аппаратный маникюр', 6000.00),
    -- Жанар Туғанбаева (77015678213) - педикюр
    ('77015678213', 'KRASOTKA_ALMATY_ALATAU', 'Аппаратный педикюр', 7500.00)
) AS data(master_phone, point_code, service_name, price)
JOIN services s ON s.point_code = data.point_code AND s.name = data.service_name;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Удаляем все реалистичные тестовые данные в правильном порядке (из-за FK constraints)
-- Сначала удаляем зависимые данные, потом родительские

-- Удаляем связи мастеров с услугами
DELETE FROM master_services WHERE updated_by = 'realistic_seed';

-- Удаляем услуги (зависят от points, service_categories, service_subcategories)
DELETE FROM services WHERE point_code IN (
    SELECT code FROM points WHERE updated_by = 'realistic_seed'
);

-- Удаляем пользователей
DELETE FROM users WHERE updated_by = 'realistic_seed';

-- Удаляем точки (зависят от networks и point_categories)
DELETE FROM points WHERE updated_by = 'realistic_seed';

-- Удаляем сети
DELETE FROM networks WHERE updated_by = 'realistic_seed';

-- Удаляем подкатегории услуг (зависят от service_categories)
DELETE FROM service_subcategories WHERE service_category_id IN (
    SELECT id FROM service_categories WHERE name IN (
        'Стрижки', 'Бритье', 'Уход за бородой', 'Укладка и стайлинг',
        'Окрашивание волос', 'Маникюр', 'Педикюр', 'Косметология', 'Массаж'
    )
);

-- Удаляем категории услуг
DELETE FROM service_categories WHERE name IN (
    'Стрижки', 'Бритье', 'Уход за бородой', 'Укладка и стайлинг',
    'Окрашивание волос', 'Маникюр', 'Педикюр', 'Косметология', 'Массаж'
);

-- Удаляем категории точек (были созданы этой миграцией)
DELETE FROM point_categories WHERE name IN ('Барбершоп', 'Салон красоты', 'Универсальный салон');

-- +goose StatementEnd
