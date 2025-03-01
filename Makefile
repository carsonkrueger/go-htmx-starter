live:
	air

all:
	make templ
	make tw

templ:
	templ generate

tw:
	npx @tailwindcss/cli -i ./public/index.css -o ./public/output.css
