package context

import (
	gctx "context"
	"time"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models"
	"github.com/carsonkrueger/main/models/auth_models"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

type PostgresTable interface {
	postgres.WritableTable
	postgres.ReadableTable
}

type PrimaryKey any

type GetBaseCols interface {
	InsertCols() postgres.ColumnList
	UpdateCols() postgres.ColumnList
	AllCols() postgres.ColumnList
	OnConflictCols() postgres.ColumnList
	UpdateOnConflictCols() []postgres.ColumnAssigment
}

type PKMatcher[PK any] interface {
	PKMatch(pk PK) postgres.BoolExpression
}

type GetUpdatedAt[R any] interface {
	GetUpdatedAt(row *R) *time.Time
}

type GetTable interface {
	Table() PostgresTable
}

type DAOBaseQueries[PK PrimaryKey, R any] interface {
	Index(ctx gctx.Context, params *models.SearchParams, db qrm.Queryable) ([]*R, error)
	GetOne(ctx gctx.Context, pk PK, db qrm.Queryable) (*R, error)
	GetMany(ctx gctx.Context, where postgres.BoolExpression, db qrm.Queryable) ([]*R, error)
	Insert(ctx gctx.Context, model *R, db qrm.Queryable) error
	InsertMany(ctx gctx.Context, models *[]*R, db qrm.Queryable) error
	Upsert(ctx gctx.Context, model *R, db qrm.Queryable) error
	UpsertMany(ctx gctx.Context, models *[]*R, db qrm.Queryable) error
	Update(ctx gctx.Context, model *R, pk PK, db qrm.Queryable) error
	Delete(ctx gctx.Context, pk PK, db qrm.Executable) error
}

type DAO[PK any, R any] interface {
	GetTable
	GetBaseCols
	PKMatcher[PK]
	GetUpdatedAt[R]
	DAOBaseQueries[PK, R]
}

type DAOManager interface {
	UsersDAO() UsersDAO
	PrivilegeDAO() PrivilegeDAO
	PrivilegeLevelsDAO() PrivilegeLevelsDAO
	SessionsDAO() SessionsDAO
	PrivilegeLevelsPrivilegesDAO() PrivilegeLevelsPrivilegesDAO
	// INSERT GET DAO
}

type UsersDAO interface {
	DAO[int64, model.Users]
	GetByEmail(email string) (*model.Users, error)
	GetPrivilegeLevelID(id int64) (*int64, error)
	GetUserPrivilegeJoinAll() (*[]auth_models.UserPrivilegeLevelJoin, error)
}

type PrivilegeDAO interface {
	DAO[int64, model.Privileges]
	GetAllJoined() ([]auth_models.JoinedPrivilegesRaw, error)
	GetPrivilegesByLevelID(levelID int64) ([]model.PrivilegeLevels, error)
}

type SessionsDAO interface {
	DAO[auth_models.SessionsPrimaryKey, model.Sessions]
}

type PrivilegeLevelsDAO interface {
	DAO[int64, model.PrivilegeLevels]
}

type PrivilegeLevelsPrivilegesDAO interface {
	DAO[auth_models.PrivilegeLevelsPrivilegesPrimaryKey, model.PrivilegeLevelsPrivileges]
}

// INSERT INTERFACE DAO
