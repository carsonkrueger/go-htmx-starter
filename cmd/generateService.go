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

func generateService() {
	service := flag.String("service", "", "camelCase Name of the service")
	flag.Parse()

	// lower first letter of table name
	service = tools.Ptr(tools.ToLowerFirst(*service))

	cfg := cfg.LoadConfig()
	lgr := logger.NewLogger(&cfg).Named("generateService")

	wd, err := os.Getwd()
	if err != nil {
		lgr.Error("failed to get working directory", zap.Error(err))
		os.Exit(1)
	}

	filePath := fmt.Sprintf("%s/services/%s.go", wd, *service)
	if err := os.MkdirAll(path.Dir(filePath), 0755); err != nil {
		lgr.Error("failed to create directory", zap.Error(err))
		os.Exit(1)
	}

	contents := ServiceFileContents(*service)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
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

func ServiceFileContents(name string) string {
	upper := tools.ToUpperFirst(name)
	lower := tools.ToLowerFirst(name)

	return fmt.Sprintf(`package services

type %[2]sService interface {
}

type %[1]s struct {
	ServiceContext
}

func New%[2]sService(ctx ServiceContext) *%[1]s {
	return &%[1]s{
		ServiceContext: ctx,
	}
}
`, lower, upper)
}
