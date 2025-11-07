package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/carsonkrueger/main/internal/templates/text"
	"github.com/carsonkrueger/main/pkg/util"
)

func generateController() {
	controller := flag.String("name", "", "camelCase Name of the controller")
	private := flag.Bool("private", true, "Is a private controller")
	flag.Parse()

	if controller == nil || *controller == "" {
		flag.Usage()
		os.Exit(1)
	}

	// lower first letter of table name
	controller = util.Ptr(util.ToLowerFirst(*controller))

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	snakeCase := util.ToSnakeCase(*controller)

	var filePath string
	if *private {
		filePath = fmt.Sprintf("%s/internal/controllers/private/%s.go", wd, snakeCase)
	} else {
		filePath = fmt.Sprintf("%s/internal/controllers/public/%s.go", wd, snakeCase)
	}
	if err := os.MkdirAll(path.Dir(filePath), 0755); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	writeControllerContents(file, *controller, *private)

	marker := ""
	newContent := ""
	if *private {
		marker = "// INSERT PRIVATE"
		newContent = fmt.Sprintf("\t\t\tprivate.New%s(appCtx),", util.ToUpperFirst(*controller))
	} else {
		marker = "// INSERT PUBLIC"
		newContent = fmt.Sprintf("\t\t\tpublic.New%s(appCtx),", util.ToUpperFirst(*controller))
	}

	appRouterPath := fmt.Sprintf("%s/internal/router/app_router.go", wd)
	util.InsertAt(appRouterPath, marker, true, newContent)
}

func writeControllerContents(f io.Writer, name string, private bool) {
	model := map[string]any{
		"Name":      util.ToUpperFirst(name),
		"NameLower": util.ToLowerFirst(name),
	}
	key := text.PrivateController
	if !private {
		key = text.PublicController
	}
	text.ExecuteTemplate(f, key, model)
}
