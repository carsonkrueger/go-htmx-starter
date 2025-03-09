package enums

import "github.com/carsonkrueger/main/gen/go_db/auth/model"

var (
	HelloWorldGet  model.Privileges = model.Privileges{Name: "helloworld/", ID: 1000}
	HelloWorldGet2 model.Privileges = model.Privileges{Name: "helloworld/get2", ID: 1001}
)
