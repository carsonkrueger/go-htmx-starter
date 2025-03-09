package services

import (
	"database/sql"
)

type IServiceManager interface {
	UsersService() IUsersService
}

type serviceManager struct {
	usersService IUsersService
}

func NewServiceManager(db *sql.DB) *serviceManager {
	return &serviceManager{
		usersService: NewUsersService(db),
	}
}

func (sm *serviceManager) UsersService() IUsersService {
	return sm.usersService
}
