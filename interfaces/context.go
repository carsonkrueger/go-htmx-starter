package interfaces

import (
	"go.uber.org/zap"
)

type INamedLogger interface {
	Lgr(name string) *zap.Logger
}

type IAppContext interface {
	INamedLogger
	SM() IServiceManager
	DM() IDAOManager
}

type ISetAppContext interface {
	SetAppCtx(ctx IAppContext)
}

type IServiceContext interface {
	INamedLogger
	DM() IDAOManager
}
