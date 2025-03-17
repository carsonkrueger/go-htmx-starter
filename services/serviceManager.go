package services

import (
	"github.com/carsonkrueger/main/interfaces"
)

type serviceManager struct {
	usersService      interfaces.IUsersService
	privilegesService interfaces.IPrivilegesService
	svcCtx            interfaces.IServiceContext
}

func NewServiceManager(svcCtx interfaces.IServiceContext) *serviceManager {
	return &serviceManager{
		svcCtx: svcCtx,
	}
}

func (sm *serviceManager) UsersService() interfaces.IUsersService {
	if sm.usersService == nil {
		sm.usersService = NewUsersService(sm.svcCtx)
	}
	return sm.usersService
}

func (sm *serviceManager) PrivilegesService() interfaces.IPrivilegesService {
	if sm.privilegesService == nil {
		sm.privilegesService = NewPrivilegesService(sm.svcCtx)
	}
	return sm.privilegesService
}
