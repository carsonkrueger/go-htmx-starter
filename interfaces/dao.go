package interfaces

import (
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/go-jet/jet/v2/postgres"
)

type IDAO interface {
	GetById(id int64) (*model.Users, error)
	Insert(row *model.Users) (int64, error)
	Upsert(row *model.Users, cols_update ...postgres.ColumnAssigment) (int64, error)
	Update(row *model.Users) error
	Delete(id int64) error
}

type IDAOManager interface {
	UsersDAO() IUsersDAO
}

type IUsersDAO interface {
	GetByEmail(email string) (*model.Users, error)
	GetById(id int64) (*model.Users, error)
	Insert(row *model.Users) (int64, error)
	Upsert(row *model.Users, colsUpdate ...postgres.ColumnAssigment) (int64, error)
	Update(row *model.Users) error
	Delete(id int64) error
	GetPrivilegeLevelID(id int64, token string) (int64, error)
	UpdateAuthToken(id int64, authToken string) error
}
