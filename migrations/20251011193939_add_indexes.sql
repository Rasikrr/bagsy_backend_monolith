-- +goose Up
-- +goose StatementBegin

-- ============================================
-- ИНДЕКСЫ (INDEXES)
-- ============================================

-- Points indexes
CREATE INDEX IF NOT EXISTS idx_points_network_code ON points(network_code);
CREATE INDEX IF NOT EXISTS idx_points_category_id ON points(category_id);

-- Users indexes
CREATE INDEX IF NOT EXISTS idx_users_point_code ON users(point_code);
CREATE INDEX IF NOT EXISTS idx_users_network_code ON users(network_code);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Services indexes
CREATE INDEX IF NOT EXISTS idx_services_point_code ON services(point_code);
CREATE INDEX IF NOT EXISTS idx_services_category_id ON services(category_id);
CREATE INDEX IF NOT EXISTS idx_services_subcategory_id ON services(subcategory_id);

-- Service subcategories indexes
CREATE INDEX IF NOT EXISTS idx_service_subcategories_service_category_id ON service_subcategories(service_category_id);

-- Master services indexes
CREATE INDEX IF NOT EXISTS idx_master_services_master_phone ON master_services(master_phone);
CREATE INDEX IF NOT EXISTS idx_master_services_service_id ON master_services(service_id);

-- Bagsies indexes
CREATE INDEX IF NOT EXISTS idx_bagsies_point_code ON bagsies(point_code);
CREATE INDEX IF NOT EXISTS idx_bagsies_client_phone ON bagsies(client_phone);
CREATE INDEX IF NOT EXISTS idx_bagsies_master_phone ON bagsies(master_phone);
CREATE INDEX IF NOT EXISTS idx_bagsies_service_id ON bagsies(service_id);
CREATE INDEX IF NOT EXISTS idx_bagsies_start_at ON bagsies(start_at);
CREATE INDEX IF NOT EXISTS idx_bagsies_end_at ON bagsies(end_at);
CREATE INDEX IF NOT EXISTS idx_bagsies_status ON bagsies(status);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop all indexes in reverse order
DROP INDEX IF EXISTS idx_bagsies_status;
DROP INDEX IF EXISTS idx_bagsies_end_at;
DROP INDEX IF EXISTS idx_bagsies_start_at;
DROP INDEX IF EXISTS idx_bagsies_service_id;
DROP INDEX IF EXISTS idx_bagsies_master_phone;
DROP INDEX IF EXISTS idx_bagsies_client_phone;
DROP INDEX IF EXISTS idx_bagsies_point_code;

DROP INDEX IF EXISTS idx_master_services_service_id;
DROP INDEX IF EXISTS idx_master_services_master_phone;

DROP INDEX IF EXISTS idx_service_subcategories_service_category_id;

DROP INDEX IF EXISTS idx_services_subcategory_id;
DROP INDEX IF EXISTS idx_services_category_id;
DROP INDEX IF EXISTS idx_services_point_code;

DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_network_code;
DROP INDEX IF EXISTS idx_users_point_code;

DROP INDEX IF EXISTS idx_points_category_id;
DROP INDEX IF EXISTS idx_points_network_code;

-- +goose StatementEnd
