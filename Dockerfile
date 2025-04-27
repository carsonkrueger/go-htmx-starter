FROM golang:1.24.0-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@latest
# DB-START
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install github.com/go-jet/jet/v2/cmd/jet@latest
RUN apk add postgresql-client
# DB-END

RUN apk add npm make
COPY package.json ./
RUN npm install

COPY . .

EXPOSE 8080

CMD ["/bin/sh", "./entrypoint.sh"]
# ENTRYPOINT ./entrypoint.sh
