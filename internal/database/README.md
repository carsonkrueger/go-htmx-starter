# database

## DAO (Data Access Object)

These objects are used to query your database.
They have a base functionality in [base_query.go](github.com/carsonkrueger/go-htmx-starter/internal/database/dao/base_query.go) that is constructed in the DAO constructor.
A DAO Manager contains all of the DAO objects for each table. [dao_manager.go](github.com/carsonkrueger/go-htmx-starter/internal/database/dao/dao_manager.go).
The DAO manager can be accessed from the AppContext/ServiceContext passed into controllers and services.
Generate a DAO automatically using `make gen-dao' and follow the instructions.
