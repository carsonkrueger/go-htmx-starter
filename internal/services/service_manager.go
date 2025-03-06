package services

import (
	"database/sql"
)

type ServiceManager struct {
	UsersService IUsersService
}

func NewServiceManager(db *sql.DB) *ServiceManager {
	return &ServiceManager{
		UsersService: NewUsersService(db),
	}
}
