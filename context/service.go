package context

import (
	"database/sql"
	"net/http"

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

type UsersService interface {
	Login(email string, password string, req *http.Request) (*string, error)
	Logout(id int64, token string) error
	LogoutRequest(req *http.Request) error
	GetAuthParts(req *http.Request) (string, int64, error)
}

// INSERT INTERFACE SERVICE
