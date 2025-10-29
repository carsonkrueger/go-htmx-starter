# database

# ⚙️ `/database` package

These files represent your **database schema**, including:
- Table definitions
- SQL builder methods
- Relationship mappings
- Query helper functions


## `/dao` subpackage (Database Access Object)

- These objects are used to query your database.

> They have a base functionality in [base_query.go](github.com/carsonkrueger/go-htmx-starter/internal/database/dao/base_query.go) that is injected in the DAO constructor.
> A DAO Manager contains all of the DAO objects for each table. [dao_manager.go](github.com/carsonkrueger/go-htmx-starter/internal/database/dao/dao_manager.go).
> The DAO manager can be accessed from the AppContext/ServiceContext passed into controllers and services.
> Generate a DAO automatically using `make generate-dao` and follow the instructions.


## `/gen` subpackage

The `/gen` directory contains **automatically generated Go code** created by [**go-jet**](https://github.com/go-jet/jet), a powerful type-safe SQL query builder for Go.

These files represent your **database schema**, including:
- Table definitions
- SQL builder methods
- Relationship mappings
- Query helper functions

> ⚠️ **Do not manually edit any files in this directory.**
>
> All files here are **regenerated automatically** whenever your database schema changes.
> Any manual changes will be **overwritten** the next time code generation runs.

---

### 🛠 gen Purpose

The generated code in `/gen` provides:
- **Type-safe SQL construction** without raw strings
- **Compile-time query validation**
- **Database schema reflection**
- **Easy model-to-table mapping**

This ensures your Go code stays consistent with your database schema, eliminating mismatches and reducing runtime SQL errors.

---

### 🧩 gen Structure

Typical layout of the `/gen` directory:
gen/
├─ db
│ ├─ schema/
│ ├─── table.go
│ ├─── users.go
│ ├─── posts.go
│ └─── ...
