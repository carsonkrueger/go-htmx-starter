package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/carsonkrueger/main/cfg"
	"github.com/carsonkrueger/main/logger"
	"github.com/carsonkrueger/main/util"
	"go.uber.org/zap"
)

func generateService() {
	service := flag.String("service", "", "camelCase Name of the service")
	flag.Parse()

	// lower first letter of table name
	service = util.Ptr(util.ToLowerFirst(*service))

	cfg := cfg.LoadConfig()
	lgr := logger.NewLogger(&cfg).Named("generateService")

	wd, err := os.Getwd()
	if err != nil {
		lgr.Error("failed to get working directory", zap.Error(err))
		os.Exit(1)
	}

	snakeCaseService := util.ToSnakeCase(*service)

	filePath := fmt.Sprintf("%s/services/%s.go", wd, snakeCaseService)
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

	upperSvcName := util.ToUpperFirst(*service) + "Service"
	svcName := util.ToLowerFirst(*service) + "Service"
	contextServicePath := fmt.Sprintf("%s/context/service.go", wd)
	serviceMgrPath := fmt.Sprintf("%s/services/service_manager.go", wd)

	util.InsertAt(contextServicePath, "// INSERT INTERFACE SERVICE", true, fmt.Sprintf("type %[1]s interface {}\n", upperSvcName))
	util.InsertAt(contextServicePath, "// INSERT GET SERVICE", true, fmt.Sprintf("\t%s() %s", upperSvcName, upperSvcName))
	util.InsertAt(serviceMgrPath, "// INSERT SERVICE", true, fmt.Sprintf("\t%s context.%s", svcName, upperSvcName))
	util.InsertAt(serviceMgrPath, "// INSERT INIT SERVICE", true, fmt.Sprintf(`func (sm *serviceManager) %[1]s() context.%[1]s {
	if sm.%[2]s == nil {
		sm.%[2]s = New%[1]s(sm.svcCtx)
	}
	return sm.%[2]s
}
`, upperSvcName, svcName))
}

func ServiceFileContents(name string) string {
	upper := util.ToUpperFirst(name)
	lower := util.ToLowerFirst(name)

	return fmt.Sprintf(`package services

import "github.com/carsonkrueger/main/context"

type %[1]sService struct {
	context.AppContext
}

func New%[2]sService(ctx context.AppContext) *%[1]sService {
	return &%[1]sService{
		ctx,
	}
}
`, lower, upper)
}
