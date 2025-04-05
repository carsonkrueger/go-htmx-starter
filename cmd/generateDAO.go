package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/carsonkrueger/main/cfg"
	"github.com/carsonkrueger/main/logger"
	"github.com/carsonkrueger/main/tools"
	"go.uber.org/zap"
)

func generateDAO() {
	table := flag.String("table", "", "camelCase Name of the DAO")
	schema := flag.String("schema", "", "camelCase Schema Name of the DAO")
	flag.Parse()

	cfg := cfg.LoadConfig()
	lgr := logger.NewLogger(&cfg).Named("generateDAO")

	if table == nil || *table == "" {
		lgr.Error("table is required")
		os.Exit(1)
	}
	if schema == nil || *schema == "" {
		lgr.Error("schema is required")
		os.Exit(1)
	}

	wd, err := os.Getwd()
	if err != nil {
		lgr.Error("failed to get working directory", zap.Error(err))
		os.Exit(1)
	}

	daoFilePath := fmt.Sprintf("%s/database/%s/%sDAO.go", wd, *schema, *table)
	contents := daoFileContents(cfg.DbConfig.Name, *schema, *table)
	file, err := os.OpenFile(daoFilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			lgr.Error("File already exists\n")
			return
		}
		lgr.Error("failed to open file", zap.Error(err))
		return
	}
	io.WriteString(file, contents)
	file.Close()
}

func daoFileContents(dbName, schema, table string) string {
	// upperSchema := tools.ToUpperFirst(schema)
	upperTable := tools.ToUpperFirst(table)

	return fmt.Sprintf(`package authDAO

import (
    "database/sql"
    "time"

    "github.com/carsonkrueger/main/gen/%[1]s/%[3]s/model"
    "github.com/carsonkrueger/main/gen/%[1]s/%[3]s/table"
    "github.com/carsonkrueger/main/tools"
    "github.com/go-jet/jet/v2/postgres"
)

type %[3]sDAO struct {
    db *sql.DB
}

func New%[3]sDAO(db *sql.DB) *%[3]sDAO {
    return &%[3]sDAO{
        db,
    }
}

func (dao *%[3]sDAO) GetById(id int64) (*model.%[3]s, error) {
    var row model.%[3]s
    err := table.%[3]s.SELECT(table.%[3]s.AllColumns).
        FROM(table.%[3]s).
        WHERE(table.%[3]s.ID.EQ(postgres.Int(id))).
        LIMIT(1).
        Query(dao.db, &row)
    if err != nil {
        return nil, err
    }
    return &row, nil
}

func (dao *%[3]sDAO) Insert(row *model.%[3]s) error {
    var res model.%[3]s
    err := table.%[3]s.
        INSERT(table.%[3]s.AllColumns.Except(table.%[3]s.ID, table.%[3]s.CreatedAt, table.%[3]s.UpdatedAt)).
        VALUES(postgres.String(row.Name)).
        RETURNING(table.%[3]s.ID).
        Query(dao.db, res)
    return err
}

func (dao *%[3]sDAO) InsertMany(rows []*model.%[3]s) error {
    if len(rows) == 0 {
        return nil
    }
    return table.%[3]s.
        INSERT(table.%[3]s.AllColumns.Except(table.%[3]s.ID, table.%[3]s.CreatedAt, table.%[3]s.UpdatedAt)).
        MODELS(rows).
        RETURNING(table.%[3]s.ID).
        Query(dao.db, &rows)
}

func (dao *%[3]sDAO) Upsert(row *model.%[3]s, colsUpdate ...postgres.ColumnAssigment) error {
    if len(colsUpdate) == 0 {
        colsUpdate = []postgres.ColumnAssigment{
            table.%[3]s.Name.SET(postgres.String(row.Name)),
        }
    }

    row.UpdatedAt = tools.Ptr(time.Now())

    return table.%[3]s.
        INSERT(table.%[3]s.AllColumns.Except(table.%[3]s.ID, table.%[3]s.CreatedAt, table.%[3]s.UpdatedAt)).
        VALUES(row.Name).
        ON_CONFLICT(table.%[3]s.Name).
        DO_UPDATE(postgres.SET(colsUpdate...)).
        RETURNING(table.%[3]s.ID).
        Query(dao.db, row)
}

func (dao *%[3]sDAO) UpsertMany(rows []*model.%[3]s, colsUpdate ...postgres.ColumnAssigment) error {
    if len(colsUpdate) == 0 {
        colsUpdate = []postgres.ColumnAssigment{
            table.%[3]s.Name.SET(table.%[3]s.Name),
        }
    }

    now := time.Now()
    for _, r := range rows {
        r.UpdatedAt = &now
    }

    return table.%[3]s.
        INSERT(table.%[3]s.AllColumns.Except(table.%[3]s.ID, table.%[3]s.CreatedAt, table.%[3]s.UpdatedAt)).
        MODELS(rows).
        ON_CONFLICT(table.%[3]s.Name).
        DO_UPDATE(postgres.SET(colsUpdate...)).
        RETURNING(table.%[3]s.ID).
        Query(dao.db, &rows)
}

func (dao *%[3]sDAO) Update(row *model.%[3]s) error {
    row.UpdatedAt = tools.Ptr(time.Now())
    _, err := table.%[3]s.
        UPDATE(table.%[3]s.EXCLUDED.ID).
        MODEL(row).
        WHERE(table.%[3]s.ID.EQ(postgres.Int(row.ID))).
        SET(table.%[3]s.UpdatedAt.SET(postgres.TimestampT(time.Now()))).
        Exec(dao.db)
    return err
}

func (dao *%[3]sDAO) Delete(id int64) error {
    _, err := table.%[3]s.DELETE().WHERE(table.%[3]s.ID.EQ(postgres.Int(id))).Exec(dao.db)
    if err != nil {
        return err
    }
    return nil
}

func (dao *%[3]sDAO) GetAll() (*[]model.%[3]s, error) {
    var rows []model.%[3]s
    err := table.%[3]s.
        SELECT(table.%[3]s.AllColumns).
        ORDER_BY(table.%[3]s.ID.DESC()).
        Query(dao.db, &rows)
    if err != nil {
        return nil, err
    }
    return &rows, nil
}
`, dbName, schema, upperTable)
}
