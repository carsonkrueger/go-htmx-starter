![webpage](https://github.com/carsonkrueger/go-htmx-starter/blob/main/webpage.png)

# ⚡ Go + HTMX Starter Kit

A powerful, starter kit for building fast and beautiful web servers using **Go**, **HTMX**, and **PostgreSQL**.
Everything you need to move from idea to production — quickly, cleanly, and with style.

---

## ✨ Core Stack

### 🐹 Go
A fast, reliable, and easy-to-learn language — perfect for modern web backends.

### 🌐 HTMX + Templ/TemplUI/Tailwind + Postgres + Jet
Build rich, reactive web interfaces using **HTML-first** development.
- **HTMX** brings interactivity without JavaScript frameworks.
- **Templ** easy type-safe UI components in Go.
- **Postgres** powerful, open-source relational database management system.
- **Jet** type-safe database access with **[go-jet](https://github.com/go-jet/jet)**, a powerful query builder that generates Go code directly from your schema.

---

## 🔒 Authentication & Authorization
Session-based authentication with role-based privilege checks.
Private routes can easily be secured with the [PrivateRouteBuilder](https://github.com/carsonkrueger/go-htmx-starter/blob/main/internal/builders/router.go), which auto-manages privilege creation and enforcement.

## PostgreSQL + Jet
Type-safe database access with **[go-jet](https://github.com/go-jet/jet)**, a powerful query builder that generates Go code directly from your schema.

## Logging with Zap
Clean, structured logging powered by **[Zap](https://github.com/uber-go/zap)**.
Logging levels are automatically configured via your `.env` environment.

## Live Reloading with Air
Enjoy fast iteration cycles using **[Air](https://github.com/air-verse/air)** — your app reloads instantly on file changes.

## Makefile Shortcuts
Automation made easy with `make` commands:
- `make live` — run the app with live reload
- `make migrate-generate` — create a new migration
- Explore more in the [Makefile »](https://github.com/carsonkrueger/go-htmx-starter/blob/main/Makefile)

# Installation
To run the application locally, follow these steps:

1. Clone the repository and navigate to the project directory.
2. Install Go and Docker.
3. Setup your `.env` using the `.env.example` file as a template.

## Running the Application with Docker

4. Build and run the Docker image using the command `make docker` - when the controllers are built and run using `make web` it will create privileges associated with each controller and insert them into the database.
5. In a separate terminal, run `make seed` to create 'basic' and 'admin' privilege levels and give all privileges to the admin level.
6. Open your browser and navigate to `http://localhost:8080` to access the application.
    - If you changed the `PORT` in your `.env` file, you will also need to update the Dockerfile.

# OR

## Running the Application Locally without Docker (Postgres still containerized)
4. Install PostgreSQL, Make, and NPM on your machine. Start your PostgreSQL service.
5. `go mod download`
6. `npm install`
7. Start the PostgreSQL container using the command: `make docker-postgres`
8. Run migrations on the postgres container: `make migrate`
9. `make jet-all` to generate all database objects and query building functionality.
10. `make live` to start the server - this will create privileges associated with each controller and insert them into the database.
11. In a separate terminal, run `make seed` to create 'basic' and 'admin' privilege levels and give all privileges to the admin level.
12. Begin live development!
