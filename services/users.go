package services

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/carsonkrueger/main/database"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/carsonkrueger/main/tools"
	"github.com/go-jet/jet/v2/postgres"
)

type IUsersService interface {
	UpdateAuthToken(id int64, authToken string) error
	Login(email string, password string) (*string, error)
	IsPermitted(userId int64, privilegeId int64) bool
}

type usersService struct {
	dm database.IDAOManager
	db *sql.DB
}

func NewUsersService(dm database.IDAOManager, db *sql.DB) *usersService {
	return &usersService{
		dm: dm,
		db: db,
	}
}

func (us *usersService) UpdateAuthToken(id int64, authToken string) error {
	_, err := table.Users.UPDATE(table.Users.AuthToken).SET(authToken).WHERE(table.Users.ID.EQ(postgres.Int(id))).Exec(us.db)
	if err != nil {
		return err
	}
	return nil
}

func (us *usersService) Login(email string, password string) (*string, error) {
	dao := us.dm.UsersDAO()
	user, err := dao.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(user.Password, "$")
	hash := tools.HashPassword(password, parts[0])

	if user.Password != hash {
		return nil, errors.New("Invalid password")
	}

	token, _ := tools.GenerateSalt()
	fullToken := fmt.Sprintf("%s$%d", token, user.ID)
	err = us.UpdateAuthToken(user.ID, fullToken)
	if err != nil {
		return nil, err
	}

	return &fullToken, nil
}

type IsPermittedResponse struct {
	PrivilegeID int64
}

func (us *usersService) IsPermitted(userId int64, privilegeId int64) bool {
	var res IsPermittedResponse

	err := table.Users.
		LEFT_JOIN(table.PrivilegeLevels, table.PrivilegeLevels.ID.EQ(table.Users.PrivilegeLevelID)).
		LEFT_JOIN(table.PrivilegeLevelsPrivileges, table.PrivilegeLevelsPrivileges.PrivilegeLevelID.EQ(table.PrivilegeLevels.ID)).
		SELECT(table.PrivilegeLevelsPrivileges.PrivilegeID.AS("IsPermittedResponse.PrivilegeID")).
		WHERE(table.PrivilegeLevelsPrivileges.PrivilegeID.EQ(postgres.Int(privilegeId)).
			AND(table.Users.ID.EQ(postgres.Int(userId)))).
		LIMIT(1).
		Query(us.db, &res)

	if err != nil {
		return false
	}

	return res.PrivilegeID == privilegeId
}
