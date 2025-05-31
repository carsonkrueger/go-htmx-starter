package context

import (
	"database/sql"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type AppContext interface {
	Lgr(name string) *zap.Logger
	SM() ServiceManager
	// DB-START
	DM() DAOManager
	DB() *sql.DB
	// DB-END
}

type appContext struct {
	Logger         *zap.Logger
	ServiceManager ServiceManager
	// DB-START
	DAOManger DAOManager
	Database  *sql.DB
	// DB-END
}

func NewAppContext(
	Logger *zap.Logger,
	ServiceManager ServiceManager,
	// DB-START
	DAOManger DAOManager,
	Database *sql.DB,
	// DB-END
) *appContext {
	return &appContext{
		Logger,
		ServiceManager,
		// DB-START
		DAOManger,
		Database,
		// DB-END
	}
}

func (ctx *appContext) Lgr(name string) *zap.Logger {
	return ctx.Logger.Named(name)
}

func (ctx *appContext) SM() ServiceManager {
	return ctx.ServiceManager
}

// DB-START
func (ctx *appContext) DM() DAOManager {
	return ctx.DAOManger
}

func (ctx *appContext) DB() *sql.DB {
	return ctx.Database
}

// DB-END

func (ctx *appContext) CleanUp() {
	if err := ctx.Logger.Sync(); err != nil {
		ctx.Logger.Error("failed to sync logger", zap.Error(err))
	}
}
