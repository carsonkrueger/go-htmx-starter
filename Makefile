# Targets with a -external suffix use the DB_EXTERNAL_PORT env variable. This is used when running the server outside of the docker container.
# Conversly targets with a -internal suffix use the DB_PORT env variable. This is used when running the server within the docker container OR if using a local database (not a docker database).

include .env

DB_URL_EXTERNAL := "postgres://${DB_USER}:${DB_PASSWORD}@${DB_EXTERNAL_HOST}:${DB_EXTERNAL_PORT}/${DB_NAME}?sslmode=disable"
DB_URL_INTERNAL := "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

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
	go build -o ./bin/main main.go

docker:
	make docker-down
	docker-compose up -d --build db go_backend --remove-orphans

docker-down:
	docker-compose down

docker-postgres:
	make docker-postgres-down
	docker-compose up -d db --remove-orphans

docker-postgres-down:
	docker-compose down db

migrate:
	migrate -database ${DB_URL_EXTERNAL} -path migrations up

migrate-internal:
	migrate -database ${DB_URL_INTERNAL} -path migrations up

migrate-down:
	migrate -database ${DB_URL_EXTERNAL} -path migrations down 1

migrate-generate:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

seed:
	go run cmd/seed.go

seed-undo:
	go run cmd/seed.go -undo=true


jet-all:
	@echo "Fetching schemas from database..."
	SCHEMAS=$$(PGPASSWORD="${DB_PASSWORD}" psql -h ${DB_EXTERNAL_HOST} -p ${DB_EXTERNAL_PORT} -U ${DB_USER} -d ${DB_NAME} -Atc "SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast')"); \
	echo "Schemas found: $$SCHEMAS"; \
	for SCHEMA in $$SCHEMAS; do \
	    echo "------ Generating models for schema: $$SCHEMA ------"; \
		jet -dsn=${DB_URL_EXTERNAL} -schema=$$SCHEMA -path=./gen; \
		make jet schema=$$SCHEMA; \
	done

jet-all-internal:
	@echo "Fetching schemas from database..."
	SCHEMAS=$$(PGPASSWORD="${DB_PASSWORD}" psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -d ${DB_NAME} -Atc "SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast')"); \
	echo "Schemas found: $$SCHEMAS"; \
	for SCHEMA in $$SCHEMAS; do \
	    echo "------ Generating models for schema: $$SCHEMA ------"; \
		jet -dsn=${DB_URL_INTERNAL} -schema=$$SCHEMA -path=./gen; \
		make jet schema=$$SCHEMA; \
	done

jet:
	jet -dsn=${DB_URL_EXTERNAL} -schema=$(schema) -path=./gen;
