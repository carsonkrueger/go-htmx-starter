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
