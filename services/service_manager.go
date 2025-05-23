package services

import (
	"database/sql"

	"github.com/carsonkrueger/main/database/daos"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models/auth_models"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"go.uber.org/zap"
)

type ServiceContext interface {
	Lgr(name string) *zap.Logger
	SM() ServiceManager
	// DB-START
	DM() daos.DAOManager
	DB() *sql.DB
	// DB-END
}

type appContext struct {
	Lgr *zap.Logger
	SM  ServiceManager
	// DB-START
	DM daos.DAOManager
	DB *sql.DB
	// DB-END
}

func NewAppContext(
	lgr *zap.Logger,
	sm ServiceManager,
	// DB-START
	dm daos.DAOManager,
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
	UserPrivilegeLevelJoinAsRowData(upl []auth_models.UserPrivilegeLevelJoin, allLevels []*model.PrivilegeLevels) []datadisplay.RowData
	JoinedPrivilegeLevelAsRowData(jpl []auth_models.JoinedPrivilegeLevel) []datadisplay.RowData
	JoinedPrivilegesAsRowData(jpl []auth_models.JoinedPrivilegesRaw) []datadisplay.RowData
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
