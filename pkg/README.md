
# 📦 `/pkg` Directory

The `/pkg` directory contains **shared, reusable packages** that can be safely imported by **external applications** or **internal modules**.
It is designed to hold logic that is **independent of the project’s core domain** and has **no dependencies on internal application code**.

> ⚠️ **Important:**
> Packages in `/pkg` should **never import from the project’s root** or internal directories (e.g. `/cmd`, `/internal`).
> They may only depend on:
> - The Go standard library
> - External third-party dependencies
> - Other packages within `/pkg`

💡 Best Practices

✅ Do:

Keep /pkg clean, generic, and portable.

Write code that could be published as its own module someday.

Keep dependencies minimal and stable.

🚫 Don’t:

Import internal project logic or configuration.

Add business rules or app-specific behavior.

Mix custom UI, API, or domain logic here.

---

## 🗂 Directory Structure

pkg/
├── db/
├── model/
├── templui/
└── util/


### `/pkg/db`
This directory contains **generated database models and query helpers**.

- Generated via tools like jet


### `/pkg/model`

This package defines project-wide models and structs used across multiple layers of the application.

- Contains types shared between backend, frontend (via API), or other custom database models.
- Should be decoupled from database or transport concerns (i.e. not ORM-specific).
- Often used for serialization, domain events, or DTOs.

### `/pkg/templui`

- Houses templui components — reusable, type-safe UI elements built with Templ

### `/pkg/util`

- Houses utility functions and helpers that are not specific to any domain or feature.
