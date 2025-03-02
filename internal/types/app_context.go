package types

import (
	"database/sql"

	"github.com/carsonkrueger/main/cfg"
	"go.uber.org/zap"
)

type AppContext struct {
	lgr *zap.Logger
	db  *sql.DB
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
	if err := ctx.lgr.Sync(); err != nil {
		ctx.lgr.Error("failed to sync logger", zap.Error(err))
	}
	if err := ctx.db.Close(); err != nil {
		ctx.lgr.Error("failed to close database", zap.Error(err))
	}
}

func (a *AppContext) GetLgr() *zap.Logger {
	return a.lgr
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
