# A Starter Kit for Web Servers Using Go + HTMX
An easy way to get started with web servers using key features: Go, HTMX, and PostgreSQL.

# Features
### Go
Go is a statically typed, compiled language designed at Google. It is a fast, efficient, and easy-to-learn language that is well-suited for building web applications.

### [HTMX](https://htmx.org/) (with [Templ](https://templ.guide)) + [Hyperscript](https://hyperscript.org/) + [TailwindCSS](https://tailwindcss.com/)
- HTMX is a library that allows you to build modern web applications using HTML, CSS, and JavaScript. It provides a set of custom attributes that you can use to add interactivity to your web pages without having to write any JavaScript code.
- Templ is a templating engine that allows you to build dynamic web pages using HTML, HTMX and Go code.
- Hyperscript is a scripting language for doing front end web development. It is designed to make it very easy to respond to events and do simple DOM manipulation in code that is directly embedded on elements on a web page.
- TailwindCSS is a utility-first CSS framework that provides a set of pre-built classes that you can use to style your web pages. It is a popular choice for building modern web applications.

### Auth with Auth Tokens and Private Route Permissions
- This project uses auth tokens in the request headers and verifies them by checking the token against the database stored in the `users` table.
- Each router endpoint with permissions attached is protected by a permission check. Users without proper permissions will be unauthorized. Attaching a permission to the
endpoint can be done using the `PrivateRouteBuilder`. See [HelloWorld](https://github.com/carsonkrueger/go-test/blob/main/internal/private_routes/hello_world2.go) private route for an example.

### Postgres
PostgreSQL is a beloved, open-source object-relational database system. It is a popular choice for building scalable and reliable web applications.

### [Jet](https://github.com/go-jet/jet)
Jet is a powerful Query Builder for Go. It provides a simple and highly customizable way to build queries and automatically generate Go code from your database.

### [Zap](https://github.com/uber-go/zap) Logging
Zap is a powerful, lightweight logging library for Go. Configured with level logging using the `APP_ENV` in your `.env` file.

### Docker
This project uses Docker to containerize the application and its dependencies. Docker provides a consistent and reproducible environment for running the application, making it easy to deploy and manage.

### [Air](https://github.com/air-verse/air)
Air enabled live reloading of the application for local development.

### Make
Make is a build automation tool that is used to automate the building of software. Many recipes are provided to help you build your application quickly and easily.
- After following the installation instructions, you can run the application using the following command: `make live`.
- Creating a new migration can be done using the following command: `make migrate-generate`.
- View the rest of the `make` targets and recipes [here](https://github.com/carsonkrueger/go-test/blob/main/Makefile).
- Targets with a -external suffix use the DB_EXTERNAL_PORT env variable. This is used when running the server outside of the docker container.
- Conversly targets with a -internal suffix use the DB_PORT env variable. This is used when running the server within the docker container OR if using a local database (not a docker database).

# Dependencies
- Go
- Docker
- Docker-Compose
- Postgres
- [Templ](https://templ.guide)
- NPM (for TailwindCSS)
- [Air](https://github.com/air-verse/air)
- [Migrate](https://github.com/golang-migrate/migrate)
- [Jet](https://github.com/go-jet/jet)

# Installation
To run the application locally, follow these steps:

1. Clone the repository and navigate to the project directory.
2. Install Go, Docker, and Docker-Compose.
3. Setup your `.env` using the `.env.example` file as a template.

## Running the Application with Docker

4. Build and run the Docker image using the command `make docker`.
5. Open your browser and navigate to `http://localhost:40080` to access the application.
    - If you changed the `DB_EXTERNAL_PORT` in your `.env` file, update this URL accordingly (`http://localhost:<DB_EXTERNAL_PORT>`)
    - If you changed the `PORT` in your `.env` file, you will also need to update the Dockerfile.

## Running the Application Locally without Docker (Postgres still containerized)
4. Install PostgreSQL, Make, and NPM on your machine. Start your PostgreSQL service.
5. `go mod download`
6. `npm install`
7. `go install github.com/a-h/templ/cmd/templ@latest`
8. `go install github.com/air-verse/air@latest`
9. `go install -tages 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
10. `go install github.com/go-jet/jet/v2/cmd/jet@latest`
11. Start the PostgreSQL container using the command: `make docker-postgres`
12. Run migrations on the postgres container: `make migrate-external`
13. `make live` to begin live development!
    - This will watch .go, .templ, and tailwind classes for changes.
14. Navigate to `http://localhost:8080` to view your changes!
    - If you changed the port in your `.env` file, update this URL accordingly.
