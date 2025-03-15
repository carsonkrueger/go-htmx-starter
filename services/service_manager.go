package services

import (
	"database/sql"

	"github.com/carsonkrueger/main/database"
)

type IServiceManager interface {
	UsersService() IUsersService
}

type serviceManager struct {
	usersService IUsersService
}

func NewServiceManager(dm database.IDAOManager, db *sql.DB) *serviceManager {
	return &serviceManager{
		usersService: NewUsersService(dm, db),
	}
}

func (sm *serviceManager) UsersService() IUsersService {
	return sm.usersService
}
