package interfaces

import (
	"database/sql"

	"go.uber.org/zap"
)

type INamedLogger interface {
	Lgr(name string) *zap.Logger
}

type IAppContext interface {
	INamedLogger
	SM() IServiceManager
	DM() IDAOManager
	DB() *sql.DB
}

type ISetAppContext interface {
	SetAppCtx(ctx IAppContext)
}

type IServiceContext interface {
	IAppContext
}
