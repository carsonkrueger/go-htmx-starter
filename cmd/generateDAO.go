package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/carsonkrueger/main/cfg"
	"github.com/carsonkrueger/main/logger"
	"github.com/carsonkrueger/main/tools"
	"go.uber.org/zap"
)

func generateDAO() {
	table := flag.String("table", "", "camelCase Name of the DAO")
	schema := flag.String("schema", "", "camelCase Schema Name of the DAO")
	flag.Parse()

	// lower first letter of table name
	table = tools.Ptr(tools.ToLowerFirst(*table))

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

	daoFilePath := fmt.Sprintf("%s/database/DAO/%s.go", wd, *table)
	if err := os.MkdirAll(path.Dir(daoFilePath), 0755); err != nil {
		lgr.Error("failed to create directory", zap.Error(err))
		os.Exit(1)
	}

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

	return fmt.Sprintf(`package DAO

import (
	"database/sql"
	"time"

	"github.com/carsonkrueger/main/gen/%[1]s/%[2]s/model"
	"github.com/carsonkrueger/main/gen/%[1]s/%[2]s/table"
	"github.com/go-jet/jet/v2/postgres"
)

type %[3]sPrimaryKey int64;

type %[4]sDAO struct {
	db *sql.DB
	DAOBaseQueries[%[3]sPrimaryKey, model.%[3]s]
}

func new%[3]sDAO(db *sql.DB) *%[4]sDAO {
	dao := &%[4]sDAO{
		db:              db,
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[%[3]sPrimaryKey, model.%[3]s](dao)
	dao.DAOBaseQueries = &queries
	return dao
}

func (dao *%[4]sDAO) Table() PostgresTable {
	return table.%[3]s
}

func (dao *%[4]sDAO) InsertCols() postgres.ColumnList {
	return table.%[3]s.AllColumns.Except(
		table.%[3]s.ID,
		table.%[3]s.CreatedAt,
		table.%[3]s.UpdatedAt,
	)
}

func (dao *%[4]sDAO) UpdateCols() postgres.ColumnList {
	return table.%[3]s.AllColumns.Except(
		table.%[3]s.ID,
		table.%[3]s.CreatedAt,
	)
}

func (dao *%[4]sDAO) AllCols() postgres.ColumnList {
	return table.%[3]s.AllColumns
}

func (dao *%[4]sDAO) OnConflictCols() postgres.ColumnList {
	return []postgres.Column{}
}

func (dao *%[4]sDAO) UpdateOnConflictCols() []postgres.ColumnAssigment {
	return []postgres.ColumnAssigment{}
}

func (dao *%[4]sDAO) PKMatch(pk %[3]sPrimaryKey) postgres.BoolExpression {
	return table.%[3]s.ID.EQ(postgres.Int64(int64(pk)))
}

func (dao *%[4]sDAO) GetUpdatedAt(row *model.%[3]s) *time.Time {
	return row.UpdatedAt
}
`, dbName, schema, upperTable, table)
}
