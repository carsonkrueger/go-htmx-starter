package interfaces

import (
	"time"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

type IPostgresTable interface {
	postgres.WritableTable
	postgres.ReadableTable
}

type PrimaryKey any

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

type IDAOBaseQueries[PK PrimaryKey, R any] interface {
	Index(params *models.SearchParams, db qrm.Queryable) ([]*R, error)
	GetOne(pk PK, db qrm.Queryable) (*R, error)
	GetMany(pk PK, db qrm.Queryable) ([]*R, error)
	Insert(model *R, db qrm.Queryable) error
	InsertMany(models *[]*R, db qrm.Queryable) error
	Upsert(model *R, db qrm.Queryable) error
	UpsertMany(models *[]*R, db qrm.Queryable) error
	Update(model *R, pk PK, db qrm.Queryable) error
	Delete(pk PK, db qrm.Executable) error
}

type IDAO[PK any, R any] interface {
	IGetTable
	IGetBaseCols
	IPKMatcher[PK]
	IGetUpdatedAt[R]
	IDAOBaseQueries[PK, R]
}

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
	GetAllJoined() ([]authModels.JoinedPrivilegesRaw, error)
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
