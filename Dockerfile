FROM golang:1.24.0-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN apk add --no-cache npm make
COPY package.json ./
RUN npm install

COPY . .

RUN make build

EXPOSE 8080

CMD ["./bin/main"]
