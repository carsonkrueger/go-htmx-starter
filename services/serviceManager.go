package services

import (
	"database/sql"

	"github.com/carsonkrueger/main/database/DAO"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"go.uber.org/zap"
)

type ServiceContext interface {
	Lgr(name string) *zap.Logger
	SM() ServiceManager
	// DB-START
	DM() DAO.DAOManager
	DB() *sql.DB
	// DB-END
}

type appContext struct {
	Lgr *zap.Logger
	SM  ServiceManager
	// DB-START
	DM DAO.DAOManager
	DB *sql.DB
	// DB-END
}

func NewAppContext(
	lgr *zap.Logger,
	sm ServiceManager,
	// DB-START
	dm DAO.DAOManager,
	db *sql.DB,
	// DB-END
) *appContext {
	return &appContext{
		lgr,
		sm,
		// DB-START
		dm,
		db,
		// DB-END
	}
}

type ServiceManager interface {
	// DB-START
	UsersService() UsersService
	PrivilegesService() PrivilegesService
	// DB-END
	// INSERT GET SERVICE
}

type PrivilegesService interface {
	CreatePrivilegeAssociation(levelID int64, privID int64) error
	DeletePrivilegeAssociation(levelID int64, privID int64) error
	CreateLevel(name string) error
	HasPermissionByID(levelID int64, permissionID int64) bool
	SetUserPrivilegeLevel(levelID int64, userID int64) error
	UserPrivilegeLevelJoinAsRowData(upl []authModels.UserPrivilegeLevelJoin, allLevels []*model.PrivilegeLevels) []datadisplay.RowData
	JoinedPrivilegeLevelAsRowData(jpl []authModels.JoinedPrivilegeLevel) []datadisplay.RowData
	JoinedPrivilegesAsRowData(jpl []authModels.JoinedPrivilegesRaw) []datadisplay.RowData
}

// INSERT INTERFACE SERVICE

type serviceManager struct {
	svcCtx ServiceContext
	// DB-START
	usersService      UsersService
	privilegesService PrivilegesService
	// DB-END
	// INSERT SERVICE
}

func NewServiceManager(svcCtx ServiceContext) *serviceManager {
	return &serviceManager{
		svcCtx: svcCtx,
	}
}

func (sm *serviceManager) SetAppContext(svcCtx ServiceContext) {
	sm.svcCtx = svcCtx
}

// DB-START
func (sm *serviceManager) UsersService() UsersService {
	if sm.usersService == nil {
		sm.usersService = NewUsersService(sm.svcCtx)
	}
	return sm.usersService
}

func (sm *serviceManager) PrivilegesService() PrivilegesService {
	if sm.privilegesService == nil {
		sm.privilegesService = NewPrivilegesService(sm.svcCtx)
	}
	return sm.privilegesService
}

// DB-END

// INSERT INIT SERVICE
