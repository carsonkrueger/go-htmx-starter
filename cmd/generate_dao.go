package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/carsonkrueger/main/internal/cfg"
	"github.com/carsonkrueger/main/internal/logger"
	"github.com/carsonkrueger/main/internal/templates/text"
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

	if table == nil || *table == "" || schema == nil || *schema == "" {
		flag.Usage()
		os.Exit(1)
	}

	wd, err := os.Getwd()
	if err != nil {
		lgr.Fatal("failed to get working directory", zap.Error(err))
	}

	snakeCaseTable := util.ToSnakeCase(*table)

	daoMgrFilePath := fmt.Sprintf("%s/internal/database/dao/dao_manager.go", wd)
	daoFilePath := fmt.Sprintf("%s/internal/database/dao/%s.go", wd, snakeCaseTable)
	if err := os.MkdirAll(path.Dir(daoFilePath), 0755); err != nil {
		lgr.Fatal("failed to create directory", zap.Error(err))
	}

	file, err := os.OpenFile(daoFilePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			lgr.Fatal("file already exists\n")
		}
		lgr.Fatal("failed to open file", zap.Error(err))
	}
	defer file.Close()
	writeDAOContents(file, cfg.DbConfig.Name, *schema, *table)

	upper := util.ToUpperFirst(*table)

	daoContextFilePath := fmt.Sprintf("%s/internal/context/dao.go", wd)
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
		dm.%[2]s = New%[1]s()
	}
	return dm.%[2]s
}
`, upperDAOName, daoName))
}

func writeDAOContents(w io.Writer, dbName, schema, table string) {
	model := map[string]string{
		"Name":      util.ToUpperFirst(table),
		"NameLower": table,
		"Schema":    schema,
		"DB":        dbName,
	}
	text.ExecuteTemplate(w, text.DAO, model)
}
