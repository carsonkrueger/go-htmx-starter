live:
	air

templ:
	templ generate

tw:
	npx @tailwindcss/cli -i ./public/index.css -o ./public/output.css

build:
	go build -o bin/main cmd/main.go

docker:
	docker-compose up -d backend
