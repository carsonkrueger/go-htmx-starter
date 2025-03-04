package types

import (
	"database/sql"

	"github.com/carsonkrueger/main/services"
)

type ServiceManager struct {
	UsersService services.IUsersService
}

func NewServiceManager(db *sql.DB) *ServiceManager {
	return &ServiceManager{
		UsersService: services.NewUsersService(db),
	}
}
