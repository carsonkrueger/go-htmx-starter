package services

import "github.com/carsonkrueger/main/internal/context"

type ServiceManager interface {
	UsersService() context.UsersService
	PrivilegesService() context.PrivilegesService
	// INSERT GET SERVICE
}

type serviceManager struct {
	*context.AppContext
	usersService      context.UsersService
	privilegesService context.PrivilegesService
	// INSERT SERVICE
}

func NewServiceManager(appCtx *context.AppContext) *serviceManager {
	return &serviceManager{
		AppContext: appCtx,
	}
}

func (sm *serviceManager) UsersService() context.UsersService {
	if sm.usersService == nil {
		sm.usersService = NewUsersService(sm.AppContext)
	}
	return sm.usersService
}

func (sm *serviceManager) PrivilegesService() context.PrivilegesService {
	if sm.privilegesService == nil {
		sm.privilegesService = NewPrivilegesService(sm.AppContext)
	}
	return sm.privilegesService
}

// INSERT INIT SERVICE
