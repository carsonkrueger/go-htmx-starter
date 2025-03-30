FROM golang:1.24.0-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go clean -cache -modcache -i -r
RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install github.com/go-jet/jet/v2/cmd/jet@latest

RUN apk add npm make postgresql-client
COPY package.json ./
RUN npm install

COPY . .

RUN make migrate
RUN make jet-all
RUN make build

EXPOSE 8080

CMD ["./bin/main"]
