.PHONY: api live_recorder test new_migration migrateup migratedown swagger docker-build docker-up docker-down

api:
	go run cmd/api/main.go

live_recorder:
	go run cmd/live_data_recorder/main.go

test:
	go test -race -v ./...

new_migration:
	@if [ -z "$(name)" ]; then \
		echo "Migration name is required. Use 'make new_migration name=your_migration_name'"; \
	else \
		migrate create -ext sql -dir internal/db/migrations -seq $(name); \
	fi

migrateup:
	docker-compose run --rm migrate up

migrateup1:
	docker-compose run --rm migrate up 1

migratedown:
	docker-compose run --rm migrate down

migratedown1:
	docker-compose run --rm migrate down 1

swagger:
	swag init -g internal/api/docs.go

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

