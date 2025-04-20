package services

import (
	"github.com/carsonkrueger/main/interfaces"
)

type serviceManager struct {
	usersService      interfaces.IUsersService
	privilegesService interfaces.IPrivilegesService
	interfaces.IAppContext
}

func NewServiceManager(appCtx interfaces.IAppContext) *serviceManager {
	return &serviceManager{
		IAppContext: appCtx,
	}
}

func (sm *serviceManager) SetAppContext(appCtx interfaces.IAppContext) {
	sm.IAppContext = appCtx
}

func (sm *serviceManager) UsersService() interfaces.IUsersService {
	if sm.usersService == nil {
		sm.usersService = NewUsersService(sm.IAppContext)
	}
	return sm.usersService
}

func (sm *serviceManager) PrivilegesService() interfaces.IPrivilegesService {
	if sm.privilegesService == nil {
		cache := NewPermissionCache()
		sm.privilegesService = NewPrivilegesService(sm.IAppContext, cache)
	}
	return sm.privilegesService
}
