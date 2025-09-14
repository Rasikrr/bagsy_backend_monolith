POSTGRES_DSN=postgres://postgres:5432@localhost:5432/bugsy?sslmode=disable
MIGRATIONS_DIR=./migrations

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" down
