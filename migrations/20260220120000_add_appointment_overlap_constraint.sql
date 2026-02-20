-- +goose Up
-- +goose StatementBegin
ALTER TABLE appointments ADD CONSTRAINT no_overlapping_appointments
EXCLUDE USING gist (
  employee_id WITH =,
  tstzrange(start_at, end_at) WITH &&
)
WHERE (status != 'cancelled');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE appointments DROP CONSTRAINT IF EXISTS no_overlapping_appointments;
-- +goose StatementEnd
