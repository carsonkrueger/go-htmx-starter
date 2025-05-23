package context

import (
	"database/sql"

	"github.com/carsonkrueger/main/database/daos"
	"github.com/carsonkrueger/main/services"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type AppContext interface {
	Lgr(name string) *zap.Logger
	SM() services.ServiceManager
	// DB-START
	DM() daos.DAOManager
	DB() *sql.DB
	// DB-END
}

type appContext struct {
	Logger         *zap.Logger
	ServiceManager services.ServiceManager
	// DB-START
	DAOManger daos.DAOManager
	Database  *sql.DB
	// DB-END
}

func NewAppContext(
	Logger *zap.Logger,
	ServiceManager services.ServiceManager,
	// DB-START
	DAOManger daos.DAOManager,
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

func (ctx *appContext) SM() services.ServiceManager {
	return ctx.ServiceManager
}

// DB-START
func (ctx *appContext) DM() daos.DAOManager {
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
