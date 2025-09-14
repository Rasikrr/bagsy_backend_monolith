-- +goose Up
-- +goose StatementBegin

CREATE TABLE points
(
    code        TEXT PRIMARY KEY,
    description TEXT,
    coordinates JSONB,
    name        TEXT NOT NULL,
    schedule    JSONB,
    created_at  TIMESTAMP DEFAULT now(),
    updated_at  TIMESTAMP DEFAULT now(),
    updated_by  TEXT
);

CREATE TABLE users
(
    phone      TEXT NOT NULL UNIQUE,
    role       TEXT NOT NULL DEFAULT 'user',
    name       TEXT NOT NULL,
    surname    TEXT,
    created_at TIMESTAMP     DEFAULT now(),
    updated_at TIMESTAMP     DEFAULT now(),
    updated_by TEXT,
    point_code TEXT
);

CREATE TABLE bagsies
(
    id         SERIAL PRIMARY KEY,
    time       TIMESTAMP NOT NULL,
    point_code TEXT      NOT NULL,
    phone TEXT      NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    updated_by TEXT
);

CREATE INDEX idx_users_point_code ON users (point_code);
CREATE INDEX idx_bagsies_point_code ON bagsies (point_code);
CREATE INDEX idx_bagsies_user_id ON bagsies (phone);


CREATE TABLE networks
(
    code        TEXT PRIMARY KEY,
    name        TEXT,
    description TEXT
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS networks;
DROP INDEX IF EXISTS idx_bagsies_user_id;
DROP INDEX IF EXISTS idx_bagsies_point_code;
DROP INDEX IF EXISTS idx_users_point_code;
DROP TABLE IF EXISTS bagsies;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS points;
-- +goose StatementEnd
