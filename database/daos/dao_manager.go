package daos

import (
	"time"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models"
	"github.com/carsonkrueger/main/models/auth_models"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"

	"database/sql"
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

type daoManager struct {
	db                            *sql.DB
	usersDAO                      UsersDAO
	privilegesDAO                 PrivilegeDAO
	privilegesLevelsDAO           PrivilegeLevelsDAO
	sessionsDAO                   SessionsDAO
	privilegesLevelsPrivilegesDAO PrivilegeLevelsPrivilegesDAO
	// INSERT DAO
}

func NewDAOManager(db *sql.DB) DAOManager {
	return &daoManager{
		db: db,
	}
}

func (dm *daoManager) UsersDAO() UsersDAO {
	if dm.usersDAO == nil {
		dm.usersDAO = NewUsersDAO(dm.db)
	}
	return dm.usersDAO
}

func (dm *daoManager) PrivilegeDAO() PrivilegeDAO {
	if dm.privilegesDAO == nil {
		dm.privilegesDAO = NewPrivilegesDAO(dm.db)
	}
	return dm.privilegesDAO
}

func (dm *daoManager) SessionsDAO() SessionsDAO {
	if dm.sessionsDAO == nil {
		dm.sessionsDAO = NewSessionsDAO(dm.db)
	}
	return dm.sessionsDAO
}

func (dm *daoManager) PrivilegeLevelsDAO() PrivilegeLevelsDAO {
	if dm.privilegesLevelsDAO == nil {
		dm.privilegesLevelsDAO = NewPrivilegeLevelsDAO(dm.db)
	}
	return dm.privilegesLevelsDAO
}

func (dm *daoManager) PrivilegeLevelsPrivilegesDAO() PrivilegeLevelsPrivilegesDAO {
	if dm.privilegesLevelsPrivilegesDAO == nil {
		dm.privilegesLevelsPrivilegesDAO = NewPrivilegeLevelsPrivilegesDAO(dm.db)
	}
	return dm.privilegesLevelsPrivilegesDAO
}

// INSERT INIT DAO
