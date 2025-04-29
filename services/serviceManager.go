package services

import (
	"database/sql"
	"net/http"

	"github.com/carsonkrueger/main/database/DAO"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
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
	UsersService() UsersService
	PrivilegesService() PrivilegesService
}

// DB-START
type UsersService interface {
	Login(email string, password string, req *http.Request) (*string, error)
	Logout(id int64, token string) error
	LogoutRequest(req *http.Request) error
	GetAuthParts(req *http.Request) (string, int64, error)
}

type PrivilegesService interface {
	BuildCache() error
	AddPermission(levelID int64, perms ...model.Privileges) error
	RemovePermission(levelID int64, perms ...model.Privileges) error
	CreateLevel(name string) error
	GetPermissions(levelID int64) []model.Privileges
	HasPermissionByID(levelID int64, permissionID int64) bool
	HasPermissionByName(levelID int64, permissionName string) bool
}

// DB-END

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
		cache := newPermissionCache()
		sm.privilegesService = NewPrivilegesService(sm.svcCtx, cache)
	}
	return sm.privilegesService
}

// DB-END
