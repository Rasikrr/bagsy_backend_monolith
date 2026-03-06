-- +goose Up
ALTER TABLE service_categories DROP COLUMN active;

-- +goose Down
ALTER TABLE service_categories ADD COLUMN active BOOLEAN DEFAULT true;
