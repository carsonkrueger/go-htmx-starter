package services

import "github.com/carsonkrueger/main/context"

type ServiceManager interface {
	// DB-START
	UsersService() context.UsersService
	PrivilegesService() context.PrivilegesService
	// DB-END
	// INSERT GET SERVICE
}

type serviceManager struct {
	svcCtx context.AppContext
	// DB-START
	usersService      context.UsersService
	privilegesService context.PrivilegesService
	// DB-END
	// INSERT SERVICE
}

func NewServiceManager(svcCtx context.ServiceContext) *serviceManager {
	return &serviceManager{
		svcCtx: svcCtx,
	}
}

func (sm *serviceManager) SetAppContext(svcCtx context.ServiceContext) {
	sm.svcCtx = svcCtx
}

// DB-START
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

// DB-END

// INSERT INIT SERVICE
