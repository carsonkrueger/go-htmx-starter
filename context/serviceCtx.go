package context

import (
	"github.com/carsonkrueger/main/interfaces"
	"go.uber.org/zap"
)

type serviceContext struct {
	lgr *zap.Logger
	dm  interfaces.IDAOManager
	pc  interfaces.IPermissionCache
}

func NewServiceContext(lgr *zap.Logger, dm interfaces.IDAOManager, pc interfaces.IPermissionCache) interfaces.IServiceContext {
	return &serviceContext{
		lgr,
		dm,
		pc,
	}
}

func (sc *serviceContext) Lgr(name string) *zap.Logger {
	return sc.lgr.Named(name)
}

func (sc *serviceContext) DM() interfaces.IDAOManager {
	return sc.dm
}

func (sc *serviceContext) PC() interfaces.IPermissionCache {
	return sc.pc
}
