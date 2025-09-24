-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE networks (
  code varchar PRIMARY KEY,
  name varchar NOT NULL,
  description text
);

CREATE TABLE points (
    code varchar PRIMARY KEY, // мб убрать primary key
    network_code varchar NOT NULL REFERENCES networks(code) ON DELETE CASCADE, // тут сильная связанность
    name varchar NOT NULL,
    description text,
    latitude double precision,
    longitude double precision,
    address varchar,
    city varchar,
    opening_hours varchar,
    schedule jsonb,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    updated_by varchar
);

CREATE INDEX idx_points_network_code ON points(network_code);
CREATE INDEX idx_points_city ON points(city);

CREATE TABLE bagsies (
     id uuid PRIMARY KEY DEFAULT gen_random_uuid(), // норм?
     time timestamp NOT NULL,
     point_code varchar NOT NULL REFERENCES points(code) ON DELETE CASCADE, // сильная связность
     phone varchar NOT NULL,
     start_at timestamp,
     end_at timestamp,
     created_at timestamp NOT NULL DEFAULT now(),
     updated_at timestamp NOT NULL DEFAULT now(),
     updated_by varchar
);

CREATE INDEX idx_bagsies_point_code ON bagsies(point_code);
CREATE INDEX idx_bagsies_phone ON bagsies(phone);
CREATE INDEX idx_bagsies_time ON bagsies(time);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS bagsies;
DROP TABLE IF EXISTS points;
DROP TABLE IF EXISTS networks;

-- +goose StatementEnd
