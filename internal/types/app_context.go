package types

import (
	"go.uber.org/zap"
)

type AppContext struct {
	lgr *zap.Logger
}

func NewAppContext(lgr *zap.Logger) *AppContext {
	return &AppContext{lgr}
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
