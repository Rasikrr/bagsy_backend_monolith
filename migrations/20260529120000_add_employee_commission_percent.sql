-- +goose Up
ALTER TABLE employees
    ADD COLUMN commission_percent INT NOT NULL DEFAULT 0
    CHECK (commission_percent BETWEEN 0 AND 100);

-- +goose Down
ALTER TABLE employees DROP COLUMN commission_percent;
