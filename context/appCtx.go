package context

import (
	"database/sql"

	"github.com/carsonkrueger/main/interfaces"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type appContext struct {
	lgr *zap.Logger
	sm  interfaces.IServiceManager
	// DB-START
	dm interfaces.IDAOManager
	db *sql.DB
	// DB-END
}

func NewAppContext(
	lgr *zap.Logger,
	sm interfaces.IServiceManager,
	// DB-START
	dm interfaces.IDAOManager,
	db *sql.DB,
	// DB-END
) *appContext {
	return &appContext{
		lgr,
		sm,
		// DB-START
		dm,
		db,
		// DB-END
	}
}

func (ctx *appContext) Lgr(name string) *zap.Logger {
	return ctx.lgr.Named(name)
}

func (ctx *appContext) SM() interfaces.IServiceManager {
	return ctx.sm
}

// DB-START
func (ctx *appContext) DM() interfaces.IDAOManager {
	return ctx.dm
}

func (ctx *appContext) DB() *sql.DB {
	return ctx.db
}

// DB-END

func (ctx *appContext) CleanUp() {
	if err := ctx.lgr.Sync(); err != nil {
		ctx.lgr.Error("failed to sync logger", zap.Error(err))
	}
}
