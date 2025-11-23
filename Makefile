include $(PWD)/.env

.PHONY: pg-migrate check-build

pg-migrate:
	@echo "Migrating PostgreSQL database"
	@migrate -path db/migrations -database $(POSTGRES_URL) up

check-build:
	@echo "Checking build"
	
	@echo "Generating SQLC code"
	@cd db && sqlc generate && cd ..

	@echo "Building db"
	@cd db && go build -o db.bin && rm db.bin && cd ..

	@echo "Building server"
	@cd server && go build -o server.bin && rm server.bin && cd ..
	
	@echo "Building producer"
	@cd producer && go build -o producer.bin && rm producer.bin && cd ..

	@echo "Building consumer"
	@cd consumer && go build -o consumer.bin && rm consumer.bin && cd ..

	@echo "Build completed"
