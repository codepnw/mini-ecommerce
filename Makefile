include dev.env

ENV_PATH=dev.env
MIGRATE_PATH=pkg/database/migrations
DATABASE_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

run:
	@go run cmd/api/main.go

docker-up:
	@docker compose --env-file=$(ENV_PATH) up -d

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
