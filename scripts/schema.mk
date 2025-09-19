MIGRATIONS_DIR=./migrations

migrate-up:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	goose -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" up

migrate-down:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	goose -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" down


create-migration:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	goose -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" create $(name) sql


