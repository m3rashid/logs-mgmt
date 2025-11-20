include $(PWD)/.env

.PHONY: pg-migrate

pg-migrate:
	@echo "Migrating PostgreSQL database"
	@migrate -path db/sql/migrations -database $(POSTGRES_URL) up
