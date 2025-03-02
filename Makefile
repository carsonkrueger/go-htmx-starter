include .env

live:
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

migrate:
	migrate -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_EXTERNAL_PORT}/${DB_NAME}?sslmode=disable" -path migrations up

migrate-down:
	migrate -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_EXTERNAL_PORT}/${DB_NAME}?sslmode=disable" -path migrations down 1

migrate-generate:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name
