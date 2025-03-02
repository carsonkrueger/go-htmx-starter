package types

import (
	"database/sql"

	"github.com/carsonkrueger/main/cfg"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type AppContext struct {
	Lgr *zap.Logger
	Db  *sql.DB
}

func NewAppContext(lgr *zap.Logger, cfg *cfg.Config) *AppContext {
	db, err := sql.Open("postgres", cfg.DbUrl())
	if err != nil {
		panic(err)
	}
	return &AppContext{
		lgr,
		db,
	}
}

func (ctx *AppContext) CleanUp() {
	if err := ctx.Lgr.Sync(); err != nil {
		ctx.Lgr.Error("failed to sync logger", zap.Error(err))
	}
	if err := ctx.Db.Close(); err != nil {
		ctx.Lgr.Error("failed to close database", zap.Error(err))
	}
}

type SetAppContext interface {
	SetAppCtx(ctx *AppContext)
}

type WithAppContext struct {
	AppCtx *AppContext
}

func (b *WithAppContext) SetAppCtx(ctx *AppContext) {
	b.AppCtx = ctx
}
