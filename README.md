# A Start Kit for Web Servers Using Go + HTMX
An easy way to get started with web servers using Go, HTMX, and PostgreSQL.

# Features
### Go
Go is a statically typed, compiled language designed at Google. It is a fast, efficient, and easy-to-learn language that is well-suited for building web applications.

### [HTMX](https://htmx.org/) (with [Templ](https://templ.guide)) + [Hyperscript](https://hyperscript.org/) + [TailwindCSS](https://tailwindcss.com/)
- HTMX is a library that allows you to build modern web applications using HTML, CSS, and JavaScript. It provides a set of custom attributes that you can use to add interactivity to your web pages without having to write any JavaScript code.
- Templ is a templating engine that allows you to build dynamic web pages using HTML, HTMX and Go code.
- Hyperscript is a scripting language for doing front end web development. It is designed to make it very easy to respond to events and do simple DOM manipulation in code that is directly embedded on elements on a web page.
- TailwindCSS is a utility-first CSS framework that provides a set of pre-built classes that you can use to style your web pages. It is a popular choice for building modern web applications.

### Auth with JWT and Permissions
- JWT (JSON Web Token) is a compact, URL-safe means of representing claims to be transferred between two parties. It is a popular choice for authentication and authorization in web applications.
- Each router endpoint with permissions attached is protected by a permission check. Users without proper permissions will be unauthorized. Attaching a permission to the
endpoint can be done using the `PrivateRouteBuilder`. See [HelloWorld](https://github.com/carsonkrueger/go-test/blob/main/internal/private_routes/hello_world2.go) private route for an example.

### Postgres
PostgreSQL is a powerful, open-source object-relational database system. It is a popular choice for building scalable and reliable web applications.

### Docker
This project uses Docker to containerize the application and its dependencies. Docker provides a consistent and reproducible environment for running the application, making it easy to deploy and manage.

### Air
Air enabled live reloading of the application for local development.

### Make
Make is a build automation tool that is used to automate the building of software. Many recipes are provided to help you build your application quickly and easily.
- After following the installation instructions, you can run the application using the following command: `make live`.
- Creating a new migration can be done using the following command: `make migrate-generate`.

# Dependencies
- Go
- Docker
- Postgres
- [Templ](https://templ.guide)
- NPM (for TailwindCSS)
- [Air](https://github.com/air-verse/air)
- [Migrate](https://github.com/golang-migrate/migrate)

# Installation
To run the application locally, follow these steps:

1. Clone the repository and navigate to the project directory.

## Running the Application with Docker

2. Setup your `.env` using the `.env.example` file as a template.
3. Build and run the Docker image using the command `make docker`.
4. Open your browser and navigate to `http://localhost:40080` to access the application.
    - If you changed the external port in your `.env` file, update this URL accordingly.

## Running the Application Locally without Docker
2. Install Go, Docker, PostgreSQL, Make, and NPM on your machine. Start your PostgreSQL service.
3. `go mod download`
4. `npm install`
5. `go install github.com/a-h/templ/cmd/templ@latest`
6. `go install github.com/air-verse/air@latest`
7. `go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
8. `make migrate`
9. `make live` to begin live development!
