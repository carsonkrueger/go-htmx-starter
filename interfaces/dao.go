package interfaces

import (
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/go-jet/jet/v2/postgres"
)

type IDAO[M any, ID any] interface {
	GetById(id ID) (*M, error)
	Insert(row *M) error
	InsertMany(rows []*M) error
	Upsert(row *M, cols_update ...postgres.ColumnAssigment) error
	UpsertMany(row []*M, cols_update ...postgres.ColumnAssigment) error
	Update(row *M) error
	Delete(id ID) error
	GetAll() (*[]M, error)
}

type IDAOManager interface {
	UsersDAO() IUsersDAO
	PrivilegeDAO() IPrivilegeDAO
	PrivilegeLevelsDAO() IPrivilegeLevelsDAO
	SessionsDAO() ISessionsDAO
	PrivilegeLevelsPrivilegesDAO() IPrivilegeLevelsPrivilegesDAO
}

type IUsersDAO interface {
	IDAO[model.Users, int64]
	GetByEmail(email string) (*model.Users, error)
	GetPrivilegeLevelID(id int64) (*int64, error)
	GetUserPrivilegeJoinAll() (*[]authModels.UserPrivilegeLevelJoin, error)
}

type IPrivilegeDAO interface {
	IDAO[model.Privileges, int64]
	GetAllJoined() (*[]authModels.JoinedPrivilegesRaw, error)
}

type ISessionsDAO interface {
	IDAO[model.Sessions, authModels.SessionsPrimaryKey]
}

type IPrivilegeLevelsDAO interface {
	IDAO[model.PrivilegeLevels, int64]
}

type IPrivilegeLevelsPrivilegesDAO interface {
	IDAO[model.PrivilegeLevelsPrivileges, authModels.PrivilegeLevelsPrivilegesPrimaryKey]
}
