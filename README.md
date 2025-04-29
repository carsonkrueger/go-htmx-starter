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

### Auth with Sessions and Private Route Privileges
- This project uses auth tokens in the request headers and verifies them by checking the token against the database stored in the `sessions` table.
- Each router endpoint with privileges attached is protected by a privilege check. Users without proper privileges will be unauthorized. Attaching a privilege to the
endpoint can be done using the `PrivateRouteBuilder`. See [UserManagement](https://github.com/carsonkrueger/go-test/blob/main/controllers/private/userManagement.go) private route for an example. Attatching a privilege to an endpoint using the `PrivateRouteBuilder` will automatically insert the privilege into the database. Each user has a privilege level. Privileges are associated with the privilege levels to know which privileges a user has. Using the `make seed` command will give all associations to the admin level.

### Postgres
PostgreSQL is a beloved, open-source object-relational database system. It is a popular choice for building scalable and reliable web applications.

### [Jet](https://github.com/go-jet/jet) Query Builder
Jet is a powerful Query Builder for Go. It provides a simple and highly customizable way to build type safe queries and automatically generate Go code from your database.

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
- Postgres
- [Templ](https://templ.guide)
- NPM (for TailwindCSS)
- [Air](https://github.com/air-verse/air)
- [Migrate](https://github.com/golang-migrate/migrate)
- [Jet](https://github.com/go-jet/jet)

# Installation
To run the application locally, follow these steps:

1. Clone the repository and navigate to the project directory.
2. Install Go and Docker.
3. Setup your `.env` using the `.env.example` file as a template.

## Running the Application with Docker

4. Build and run the Docker image using the command `make docker` - when the controllers are built and run using `make web` it will create privileges associated with each controller and insert them into the database.
5. In a separate terminal, run `make seed` to create 'basic' and 'admin' privilege levels and give all privileges to the admin level.
6. Open your browser and navigate to `http://localhost:40080` to access the application.
    - If you changed the `DB_EXTERNAL_PORT` in your `.env` file, update this URL accordingly (`http://localhost:<DB_EXTERNAL_PORT>`)
    - If you changed the `PORT` in your `.env` file, you will also need to update the Dockerfile.

# OR

## Running the Application Locally without Docker (Postgres still containerized)
4. Install PostgreSQL, Make, and NPM on your machine. Start your PostgreSQL service.
5. `go mod download`
6. `npm install`
7. Install system dependencies with `make install-system-deps`
8. Start the PostgreSQL container using the command: `make docker-postgres`
9. Run migrations on the postgres container: `make migrate`
10. `make jet-all` to generate all database objects and query building functionality.
11. `make live` to start the server - this will create privileges associated with each controller and insert them into the database.
12. In a separate terminal, run `make seed` to create 'basic' and 'admin' privilege levels and give all privileges to the admin level.
13. Begin live development!

# Don't want to use a database?

### No Problem!
Run `make remove-db-files` to remove all files and source related to the database.
This uses markers to define areas in the source code that are database specific.
You will still need to remove all the imports that are no longer needed/existing.

# Don't like the 'DB-START' markers in your source code?

### No Problem!
Run `make remove-markers` to remove all markers.
