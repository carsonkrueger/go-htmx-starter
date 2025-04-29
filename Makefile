# Targets with a -external suffix use the DB_EXTERNAL_PORT env variable. This is used when running the server outside of the docker container.
# Conversly targets with a -internal suffix use the DB_PORT env variable. This is used when running the server within the docker container OR if using a local database (not a docker database).

include .env

GO_BIN_PATH := ~/go/bin
AIR_CMD := ${GO_BIN_PATH}/air
TEMPL_CMD := ${GO_BIN_PATH}/templ

# DB-START
MIGRATE_CMD := ${GO_BIN_PATH}/migrate
JET_CMD := ${GO_BIN_PATH}/jet
DB_URL_EXTERNAL := "postgres://${DB_USER}:${DB_PASSWORD}@${DB_EXTERNAL_HOST}:${DB_EXTERNAL_PORT}/${DB_NAME}?sslmode=disable"
DB_URL_INTERNAL := "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
# DB-END

live:
# DB-START
	make docker-postgres
# DB-END
	${AIR_CMD}

templ:
	${TEMPL_CMD} generate

tw:
	npx @tailwindcss/cli -i ./public/index.css -o ./public/output.css

build:
	make tw
	make templ
	go build -o ./bin/main .

docker:
	make docker-down
	docker compose up -d --build \
# DB-START
	    db \
# DB-END
		go_backend \
		--remove-orphans

docker-down:
	docker compose down

# DB-START
docker-postgres:
	make docker-postgres-down
	docker compose up -d db --remove-orphans

docker-postgres-down:
	docker compose down db

migrate:
	${MIGRATE_CMD} -database ${DB_URL_EXTERNAL} -path migrations up

migrate-internal:
	${MIGRATE_CMD} -database ${DB_URL_INTERNAL} -path migrations up

migrate-down:
	${MIGRATE_CMD} -database ${DB_URL_EXTERNAL} -path migrations down 1

migrate-generate:
	@read -p "Enter migration name: " name; \
	${MIGRATE_CMD} create -ext sql -dir migrations -seq $$name

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
# DB-END

generate-public-controller:
	@echo "Enter camelCase controller name: "; \
	read controller; \
	go run . -name="$$controller" -private=false genController

# DB-START
seed:
	go run . seed

seed-undo:
	go run . -undo=true seed

jet-all:
	@echo "Fetching schemas from database..."
	SCHEMAS=$$(PGPASSWORD="${DB_PASSWORD}" psql -h ${DB_EXTERNAL_HOST} -p ${DB_EXTERNAL_PORT} -U ${DB_USER} -d ${DB_NAME} -Atc "SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast')"); \
	echo "Schemas found: $$SCHEMAS"; \
	for SCHEMA in $$SCHEMAS; do \
	    echo "------ Generating models for schema: $$SCHEMA ------"; \
		${JET_CMD} -dsn=${DB_URL_EXTERNAL} -schema=$$SCHEMA -path=./gen; \
	done

jet-all-internal:
	@echo "Fetching schemas from database..."
	SCHEMAS=$$(PGPASSWORD="${DB_PASSWORD}" psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -d ${DB_NAME} -Atc "SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('pg_catalog', 'information_schema', 'pg_toast')"); \
	echo "Schemas found: $$SCHEMAS"; \
	for SCHEMA in $$SCHEMAS; do \
	    echo "------ Generating models for schema: $$SCHEMA ------"; \
		${JET_CMD} -dsn=${DB_URL_INTERNAL} -schema=$$SCHEMA -path=./gen; \
	done

jet:
	${JET_CMD} -dsn=${DB_URL_EXTERNAL} -schema=$(schema) -path=./gen;
# DB-END

install-system-deps:
	go install github.com/a-h/templ/cmd/templ@latest` \
	go install github.com/air-verse/air@latest` \
	# DB-START
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest` \
	go install github.com/go-jet/jet/v2/cmd/jet@latest`
	# DB-END

remove-markers:
	@find . \( $(foreach dir,$(EXCLUDE_DIRS),-path ./$(dir) -o ) -false \) -prune -o -type f -exec sed -i '/[\/#]\s*DB-START\s*$$/d; /[\/#]\s*DB-END\s*$$/d' {} +

EXCLUDE_DIRS=volumes node_modules

remove-db-files:
	rm -f ./builders/handler.go
	rm -f ./context/context.go
	rm -f ./cmd/generateDAO.go
	rm -f ./cmd/seed.go
	rm -rf ./controllers/private
	rm -rf ./controllers/public/login.go
	rm -rf ./controllers/public/signup.go
	rm -rf ./database
	rm -rf ./gen
	rm -rf ./interfaces/dao.go
	rm -rf ./middlewares
	rm -rf ./migrations
	rm -rf ./models/authModels
	rm -f ./models/database.go
	rm -rf ./seeders
	rm -rf ./constant
	rm -f ./services/privileges.go
	rm -f ./services/users.go
	make cut-db-start-end

cut-db-start-end:
	@find . \( $(foreach dir,$(EXCLUDE_DIRS),-path ./$(dir) -o ) -false \) -prune -o -type f -exec sed -i '/[\/#]\s*DB-START/,/[\/#]\s*DB-END/d' {} +
