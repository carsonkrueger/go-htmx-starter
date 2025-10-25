package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/carsonkrueger/main/internal/gen/template"
	"github.com/carsonkrueger/main/pkg/util"
)

func generateController() {
	controller := flag.String("name", "", "camelCase Name of the controller")
	private := flag.Bool("private", true, "Is a private controller")
	flag.Parse()

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

	writeContents(file, *controller, *private)

	upper := util.ToUpperFirst(*controller)

	marker := ""
	newContent := ""
	if *private {
		marker = "// INSERT PRIVATE"
		newContent = fmt.Sprintf("\t\t\tprivate.New%s(appCtx),", upper)
	} else {
		marker = "// INSERT PUBLIC"
		newContent = fmt.Sprintf("\t\t\tpublic.New%s(appCtx),", upper)
	}

	appRouterPath := fmt.Sprintf("%s/router/app_router.go", wd)
	util.InsertAt(appRouterPath, marker, true, newContent)
}

type ControllerModel struct {
	Name      string
	NameLower string
}

func writeContents(f io.Writer, name string, private bool) {
	model := ControllerModel{
		Name:      util.ToUpperFirst(name),
		NameLower: util.ToLowerFirst(name),
	}
	key := template.PrivateController
	if !private {
		key = template.PublicController
	}
	fmt.Printf("Writing contents for %s controller: %v %s\n", model.Name, model, key)
	template.ExecuteTemplate(f, key, model)
}
