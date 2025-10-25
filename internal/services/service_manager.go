package services

import "github.com/carsonkrueger/main/internal/context"

type ServiceManager interface {
	UsersService() context.UsersService
	PrivilegesService() context.PrivilegesService
	// INSERT GET SERVICE
}

type serviceManager struct {
	svcCtx            context.AppContext
	usersService      context.UsersService
	privilegesService context.PrivilegesService
	// INSERT SERVICE
}

func NewServiceManager(svcCtx context.AppContext) *serviceManager {
	return &serviceManager{
		svcCtx: svcCtx,
	}
}

func (sm *serviceManager) SetAppContext(svcCtx context.AppContext) {
	sm.svcCtx = svcCtx
}

func (sm *serviceManager) UsersService() context.UsersService {
	if sm.usersService == nil {
		sm.usersService = NewUsersService(sm.svcCtx)
	}
	return sm.usersService
}

func (sm *serviceManager) PrivilegesService() context.PrivilegesService {
	if sm.privilegesService == nil {
		sm.privilegesService = NewPrivilegesService(sm.svcCtx)
	}
	return sm.privilegesService
}

// INSERT INIT SERVICE
