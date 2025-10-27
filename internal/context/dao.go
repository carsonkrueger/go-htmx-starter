package context

import (
	gctx "context"
	"time"

	"github.com/carsonkrueger/main/internal/constant"
	"github.com/carsonkrueger/main/pkg/db/auth/model"
	model1 "github.com/carsonkrueger/main/pkg/model"
	"github.com/go-jet/jet/v2/postgres"
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
	GetAll(ctx gctx.Context) ([]R, error)
	GetOne(ctx gctx.Context, pk PK) (R, error)
	GetMany(ctx gctx.Context, pks []PK) ([]R, error)
	Insert(ctx gctx.Context, model *R) error
	InsertMany(ctx gctx.Context, models []R) error
	Upsert(ctx gctx.Context, model *R) error
	UpsertMany(ctx gctx.Context, models []R) error
	Update(ctx gctx.Context, model *R, pk PK) error
	Delete(ctx gctx.Context, pk PK) error
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
	RolesDAO() RolesDAO
	SessionsDAO() SessionsDAO
	RolesPrivilegesDAO() RolesPrivilegesDAO
	// INSERT GET DAO
}

type UsersDAO interface {
	DAO[int64, model.Users]
	GetByEmail(ctx gctx.Context, email string) (*model.Users, error)
	GetRoleID(ctx gctx.Context, id int64) (*int64, error)
	GetUserPrivilegeJoinAll(ctx gctx.Context) (*[]model1.UserRoleJoin, error)
}

type PrivilegeDAO interface {
	DAO[int64, model.Privileges]
	GetAllJoined(ctx gctx.Context) ([]model1.JoinedPrivilegesRaw, error)
	GetPrivilegesByRoleID(ctx gctx.Context, roleID int64) ([]model.Roles, error)
	GetManyByName(ctx gctx.Context, names []constant.PrivilegeName) ([]model.Privileges, error)
}

type SessionsDAO interface {
	DAO[model1.SessionsPrimaryKey, model.Sessions]
}

type RolesDAO interface {
	DAO[int16, model.Roles]
}

type RolesPrivilegesDAO interface {
	DAO[model1.RolesPrivilegesPrimaryKey, model.RolesPrivileges]
}

// INSERT INTERFACE DAO
