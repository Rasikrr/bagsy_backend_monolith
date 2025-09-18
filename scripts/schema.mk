MIGRATIONS_DIR=./migrations

migrate-up:
	echo $(POSTGRES_DSN)
	goose -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" down


create-migration:
	goose -dir $(MIGRATIONS_DIR) postgres "$(POSTGRES_DSN)" create $(name) sql


