package interfaces

import (
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/go-jet/jet/v2/postgres"
)

type IDAO[M any] interface {
	GetById(id int64) (*M, error)
	Insert(row *M) (int64, error)
	Upsert(row *M, cols_update ...postgres.ColumnAssigment) (int64, error)
	Update(row *M) error
	Delete(id int64) error
	GetAll() ([]*M, error)
}

type IDAOManager interface {
	UsersDAO() IUsersDAO
	PrivilegeDAO() IPrivilegeDAO
}

type IUsersDAO interface {
	IDAO[model.Users]
	GetByEmail(email string) (*model.Users, error)
	GetPrivilegeLevelID(id int64, token string) (int64, error)
	UpdateAuthToken(id int64, authToken string) error
}

type IPrivilegeDAO interface {
	IDAO[model.Privileges]
	GetAllJoined() ([]authModels.JoinedPrivilegesRaw, error)
}
