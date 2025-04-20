package interfaces

import (
	"time"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/go-jet/jet/v2/postgres"
)

type IPostgresTable interface {
	postgres.WritableTable
	postgres.ReadableTable
}

type PK any

type IGetBaseCols interface {
	InsertCols() postgres.ColumnList
	UpdateCols() postgres.ColumnList
	AllCols() postgres.ColumnList
	OnConflictCols() postgres.ColumnList
	UpdateOnConflictCols() []postgres.ColumnAssigment
}

type IPKMatcher[PK any] interface {
	PKMatch(pk PK) postgres.BoolExpression
}

type IGetUpdatedAt[R any] interface {
	GetUpdatedAt(row *R) *time.Time
}

type IGetTable interface {
	Table() IPostgresTable
}

type IDAO[PK any, R any] interface {
	IGetTable
	IGetBaseCols
	IPKMatcher[PK]
	IGetUpdatedAt[R]
}

// type IDAO[M any, ID any] interface {
// GetById(id ID) (*M, error)
// Insert(row *M) error
// InsertMany(rows []*M) error
// Upsert(row *M, cols_update ...postgres.ColumnAssigment) error
// UpsertMany(row []*M, cols_update ...postgres.ColumnAssigment) error
// Update(row *M) error
// Delete(id ID) error
// GetAll() (*[]M, error)
// }

type IDAOManager interface {
	UsersDAO() IUsersDAO
	PrivilegeDAO() IPrivilegeDAO
	PrivilegeLevelsDAO() IPrivilegeLevelsDAO
	SessionsDAO() ISessionsDAO
	PrivilegeLevelsPrivilegesDAO() IPrivilegeLevelsPrivilegesDAO
}

type IUsersDAO interface {
	IDAO[int64, model.Users]
	GetByEmail(email string) (*model.Users, error)
	GetPrivilegeLevelID(id int64) (*int64, error)
	GetUserPrivilegeJoinAll() (*[]authModels.UserPrivilegeLevelJoin, error)
}

type IPrivilegeDAO interface {
	IDAO[int64, model.Privileges]
	GetAllJoined() (*[]authModels.JoinedPrivilegesRaw, error)
}

type ISessionsDAO interface {
	IDAO[authModels.SessionsPrimaryKey, model.Sessions]
}

type IPrivilegeLevelsDAO interface {
	IDAO[int64, model.PrivilegeLevels]
}

type IPrivilegeLevelsPrivilegesDAO interface {
	IDAO[authModels.PrivilegeLevelsPrivilegesPrimaryKey, model.PrivilegeLevelsPrivileges]
}
