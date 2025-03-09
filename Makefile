# Targets with a -external suffix use the DB_EXTERNAL_PORT env variable. This is used when running the server outside of the docker container.
# Conversly targets with a -internal suffix use the DB_PORT env variable. This is used when running the server within the docker container OR if using a local database (not a docker database).

include .env

live:
	make docker-postgres
	air

templ:
	templ generate

tw:
	npx @tailwindcss/cli -i ./public/index.css -o ./public/output.css

build:
	make tw
	make templ
	go build -o bin/main cmd/main.go

docker:
	docker-compose up -d go_backend go_db --remove-orphans

docker-postgres:
	docker-compose up -d go_db --remove-orphans

migrate-external migrate:
	migrate -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_EXTERNAL_PORT}/${DB_NAME}?sslmode=disable" -path migrations up

migrate-internal:
	migrate -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -path migrations up

migrate-down-external migrate-down:
	migrate -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_EXTERNAL_PORT}/${DB_NAME}?sslmode=disable" -path migrations down 1

migrate-down-internal:
	migrate -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -path migrations down 1

migrate-generate:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

seed-external seed:
	go run cmd/seeds/seed.go -local=true

seed-internal:
	go run cmd/seeds/seed.go

seed-undo-external seed-undo:
	go run cmd/seeds/seed.go -local=true -undo=true

seed-undo-internal:
	go run cmd/seeds/seed.go -undo=true


jet-all-external jet-all:
	@echo "Fetching schemas from database..."
	@DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_EXTERNAL_PORT}/${DB_NAME}?sslmode=disable"; \
	SCHEMAS=$$(PGPASSWORD="${DB_PASSWORD}" psql -h ${HOST} -p ${DB_EXTERNAL_PORT} -U ${DB_USER} -d ${DB_NAME} -Atc "SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast')"); \
	echo "Schemas found: $$SCHEMAS"; \
	for SCHEMA in $$SCHEMAS; do \
	    echo "------ Generating models for schema: $$SCHEMA ------"; \
		jet -dsn="$$DB_URL" -schema=$$SCHEMA -path=./gen; \
		make jet schema=$$SCHEMA; \
	done

jet-all-internal:
	@echo "Fetching schemas from database..."
	@DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"; \
	SCHEMAS=$$(PGPASSWORD="${DB_PASSWORD}" psql -h ${HOST} -p ${DB_EXTERNAL_PORT} -U ${DB_USER} -d ${DB_NAME} -Atc "SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast')"); \
	echo "Schemas found: $$SCHEMAS"; \
	for SCHEMA in $$SCHEMAS; do \
		echo "------ Generating models for schema: $$SCHEMA ------"; \
		make jet-internal schema=$$SCHEMA; \
	done

jet-external jet:
	@DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_EXTERNAL_PORT}/${DB_NAME}?sslmode=disable"; \
	jet -dsn="$$DB_URL" -schema=$(schema) -path=./gen;

jet-internal:
	@DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"; \
	jet -dsn="$$DB_URL" -schema=$(schema) -path=./gen;
