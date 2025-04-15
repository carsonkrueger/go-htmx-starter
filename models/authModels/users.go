package authModels

import (
	"fmt"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models"
	"github.com/carsonkrueger/main/templates/datadisplay"
)

type UserPrivilegeLevelJoin struct {
	model.Users
	PLID   int64
	PLName string
}

func (up *UserPrivilegeLevelJoin) AsCellData() []datadisplay.CellData {
	return []datadisplay.CellData{
		{
			ID:    fmt.Sprintf("%d-%d-name", up.Users.ID, up.PLID),
			Width: 1,
			Body:  datadisplay.Text(fmt.Sprintf("%s, %s", up.Users.LastName, up.Users.FirstName), models.MD),
		},
		{
			ID:    fmt.Sprintf("%d-%d-pl", up.Users.ID, up.PLID),
			Width: 1,
			Body:  datadisplay.Text(up.PLName, models.MD),
		},
		{
			ID:    fmt.Sprintf("%d-%d-ca", up.Users.ID, up.PLID),
			Width: 1,
			Body:  datadisplay.Text(up.Users.CreatedAt.GoString(), models.MD),
		},
	}
}
