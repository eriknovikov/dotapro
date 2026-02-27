

MIGRATIONS_PATH = ./migrations

.PHONY: migrate-up migrate-down migrate-create migrate-force

# Run all pending migrations
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database $(LOCAL_DB_URL) up

# Rollback the last migration
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database $(LOCAL_DB_URL) down 2

# Create a new migration file pair (usage: make migrate-new name=add_users_table)
migrate-new:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -digits 3 -seq $(name)

# Fix a "dirty" database state (usage: make migrate-force v=1)
migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database $(LOCAL_DB_URL) force $(v)

# Check current migration version
migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database $(LOCAL_DB_URL) version