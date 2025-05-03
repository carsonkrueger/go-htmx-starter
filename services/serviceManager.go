package services

import (
	"database/sql"

	"github.com/carsonkrueger/main/database/DAO"
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
}

type serviceManager struct {
	// DB-START
	usersService      UsersService
	privilegesService PrivilegesService
	// DB-END
	svcCtx ServiceContext
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
