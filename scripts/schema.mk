MIGRATIONS_DIR := ./migrations
GOOSE := go run github.com/pressly/goose/v3/cmd/goose@v3.25.0

# Создать файл миграции
migration:
	$(GOOSE) -dir $(MIGRATIONS_DIR) create $(name) sql

migrate-up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" up

# Миграция до определенной версии (например, VERSION=20240101)
migrate-to:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" up-to $(VERSION)

# Откат до конкретной версии
migrate-down-to:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" down-to $(VERSION)

# Откат всех миграций в ноль
migrate-reset:
	$(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" reset


