# Targets with a -external suffix use the DB_EXTERNAL_PORT env variable. This is used when running the server outside of the docker container.
# Conversly targets with a -internal suffix use the DB_PORT env variable. This is used when running the server within the docker container OR if using a local database (not a docker database).

include .env

GO_BIN_PATH := ~/go/bin
AIR_CMD := ${GO_BIN_PATH}/air
TEMPL_CMD := ${GO_BIN_PATH}/templ

MIGRATE_CMD := ${GO_BIN_PATH}/migrate
JET_CMD := ${GO_BIN_PATH}/jet
DB_URL_EXTERNAL := "postgres://${DB_USER}:${DB_PASSWORD}@${DB_EXTERNAL_HOST}:${DB_EXTERNAL_PORT}/${DB_NAME}?sslmode=disable"
DB_URL_INTERNAL := "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
JET_MODEL_PATH := ../../../../pkg/model/db

live:
	make docker-postgres
	${AIR_CMD}

templ:
	${TEMPL_CMD} generate

tw:
	npx @tailwindcss/cli -i ./public/index.css -o ./public/output.css

build:
	make tw
	make templ
	go build -ldflags="-s -w" -o ./bin/main .

docker:
	make docker-down
	docker compose up -d --build \
	    go_starter_db \
		go_starter_backend \
		--remove-orphans

docker-down:
	docker compose down

docker-postgres:
	make docker-postgres-down
	docker compose up -d go_starter_db --remove-orphans

docker-postgres-down:
	docker compose down go_starter_db

migrate:
	${MIGRATE_CMD} -database ${DB_URL_EXTERNAL} -path migrations up

migrate-internal:
	${MIGRATE_CMD} -database ${DB_URL_INTERNAL} -path migrations up

migrate-down:
	${MIGRATE_CMD} -database ${DB_URL_EXTERNAL} -path migrations down 1

migrate-generate:
	@read -p "Enter migration name: " name; \
	${MIGRATE_CMD} create -ext sql -dir migrations -seq $$name

generate-service:
	@echo "Enter camelCase service name: "; \
	read service; \
	go run . -service="$$service" genService

generate-dao:
	@echo "Enter camelCase table name: "; \
	read table; \
	echo "Enter schema name: "; \
	read schema; \
	go run . -schema="$$schema" -table="$$table" genDAO

generate-private-controller:
	@echo "Enter camelCase controller name: "; \
	read controller; \
	go run . -name="$$controller" -private=true genController

generate-public-controller:
	@echo "Enter camelCase controller name: "; \
	read controller; \
	go run . -name="$$controller" -private=false genController

# seed:
	go run . seed

seed-undo:
	go run . -undo=true seed

jet-all:
	@echo "Fetching schemas from database..."
	SCHEMAS=$$(PGPASSWORD='${DB_PASSWORD}' psql -h ${DB_EXTERNAL_HOST} -p ${DB_EXTERNAL_PORT} -U ${DB_USER} -d ${DB_NAME} -Atc "SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast')"); \
	echo "Schemas found: $$SCHEMAS"; \
	for SCHEMA in $$SCHEMAS; do \
	    echo "------ Generating models for schema: $$SCHEMA ------"; \
		${JET_CMD} -dsn=${DB_URL_EXTERNAL} -schema=$$SCHEMA -rel-model-path=${JET_MODEL_PATH}/$$SCHEMA -path=./internal/gen; \
	done

jet-all-internal:
	@echo "Fetching schemas from database..."
	SCHEMAS=$$(PGPASSWORD='${DB_PASSWORD}' psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -d ${DB_NAME} -Atc "SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast')"); \
	echo "Schemas found: $$SCHEMAS"; \
	for SCHEMA in $$SCHEMAS; do \
	    echo "------ Generating models for schema: $$SCHEMA ------"; \
		${JET_CMD} -dsn=${DB_URL_INTERNAL} -schema=$$SCHEMA -rel-model-path=${JET_MODEL_PATH}/$$SCHEMA -path=./internal/gen; \
	done
