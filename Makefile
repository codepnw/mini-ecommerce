include .env

ENV_PATH=dev.env
MIGRATE_PATH=pkg/database/migrations
DATABASE_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

run:
	@go run cmd/api/main.go

test-usecase:
	@echo "Running ONLY Usecase tests..."
	go test ./internal/*/usecase -cover

# NOTE: -n check if not exists
cp-env:
	@cp -n .env.example .env

# Delete Volumn
docker-down:
	@docker-compose down -v

docker-up:
	@docker compose up --build -d

migrate-create:
	@migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

migrate-up:
	@migrate -database $(DATABASE_URL) -path $(MIGRATE_PATH) up

migrate-down:
	@migrate -database $(DATABASE_URL) -path $(MIGRATE_PATH) down 1

migrate-force:
	@migrate -database $(DATABASE_URL) -path $(MIGRATE_PATH) force $(version)

migrate-reset:
	@migrate -database $(DATABASE_URL) -path $(MIGRATE_PATH) down
