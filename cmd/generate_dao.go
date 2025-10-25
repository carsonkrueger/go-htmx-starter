package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/carsonkrueger/main/internal/cfg"
	"github.com/carsonkrueger/main/internal/logger"
	"github.com/carsonkrueger/main/pkg/util"
	"go.uber.org/zap"
)

func generateDAO() {
	table := flag.String("table", "", "camelCase name of the table")
	schema := flag.String("schema", "", "camelCase name of the schema name")
	flag.Parse()

	// lower first letter of table name
	table = util.Ptr(util.ToLowerFirst(*table))

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

	snakeCaseTable := util.ToSnakeCase(*table)

	daoMgrFilePath := fmt.Sprintf("%s/database/dao/dao_manager.go", wd)
	daoFilePath := fmt.Sprintf("%s/database/dao/%s.go", wd, snakeCaseTable)
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

	upper := util.ToUpperFirst(*table)

	daoContextFilePath := fmt.Sprintf("%s/context/dao.go", wd)
	upperDAOName := upper + "DAO"
	daoName := *table + "DAO"

	util.InsertAt(daoContextFilePath, "// INSERT GET DAO", true, fmt.Sprintf("\t%s() %s", upperDAOName, upperDAOName))
	util.InsertAt(daoContextFilePath, "// INSERT INTERFACE DAO", true, fmt.Sprintf(`type %[1]s interface {
	DAO[int64, model.%[2]s]
}
`, upperDAOName, upper))
	util.InsertAt(daoMgrFilePath, "// INSERT DAO", true, fmt.Sprintf("\t%s context.%s", daoName, upperDAOName))
	util.InsertAt(daoMgrFilePath, "// INSERT INIT DAO", true, fmt.Sprintf(`func (dm *daoManager) %[1]s() context.%[1]s {
	if dm.%[2]s == nil {
		dm.%[2]s = New%[1]s(dm.db)
	}
	return dm.%[2]s
}
`, upperDAOName, daoName))
}

func daoFileContents(dbName, schema, table string) string {
	// upperSchema := util.ToUpperFirst(schema)
	upperTable := util.ToUpperFirst(table)

	return fmt.Sprintf(`package dao

import (
	"database/sql"
	"time"

	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/gen/%[1]s/%[2]s/model"
	"github.com/carsonkrueger/main/gen/%[1]s/%[2]s/table"
	"github.com/go-jet/jet/v2/postgres"
)

type %[3]sPrimaryKey int64;

type %[4]sDAO struct {
	db *sql.DB
	context.DAOBaseQueries[%[3]sPrimaryKey, model.%[3]s]
}

func New%[3]sDAO(db *sql.DB) *%[4]sDAO {
	dao := &%[4]sDAO{
		db:              db,
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[%[3]sPrimaryKey, model.%[3]s](dao)
	dao.DAOBaseQueries = &queries
	return dao
}

func (dao *%[4]sDAO) Table() context.PostgresTable {
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
