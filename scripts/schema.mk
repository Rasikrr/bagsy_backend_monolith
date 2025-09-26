MIGRATIONS_DIR := ./migrati.ons
GOOSE := go run github.com/pressly/goose/v3/cmd/goose@v3.25.0

migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" up

migrate-down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" down

create-migration:
	$(GOOSE) -dir $(MIGRATIONS_DIR) create $(name) sql
