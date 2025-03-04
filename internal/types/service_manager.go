package types

import (
	"database/sql"

	"github.com/carsonkrueger/main/services"
)

type IServiceManager interface {
	UsersService() services.IUsersService
}

type serviceManager struct {
	usersService services.IUsersService
}

func NewServiceManager(db *sql.DB) IServiceManager {
	return &serviceManager{
		usersService: services.NewUsersService(db),
	}
}

func (sm *serviceManager) UsersService() services.IUsersService {
	return sm.usersService
}
