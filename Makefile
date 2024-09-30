.PHONY: new_migration migrateup migratedown migrateup1 migratedown1

new_migration:
	@if [ -z "$(name)" ]; then \
		echo "Migration name is required. Use 'make new_migration name=your_migration_name'"; \
	else \
		migrate create -ext sql -dir internal/db/migrations -seq $(name); \
	fi

migrateup:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose down 1

