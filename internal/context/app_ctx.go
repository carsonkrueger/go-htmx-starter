package context

import (
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type AppContext struct {
	logger         *zap.Logger
	ServiceManager ServiceManager
	DAOManger      DAOManager
}

func NewAppContext(
	Logger *zap.Logger,
	ServiceManager ServiceManager,
	DAOManger DAOManager,
) *AppContext {
	return &AppContext{
		Logger,
		ServiceManager,
		DAOManger,
	}
}

func (ctx *AppContext) Lgr(name string) *zap.Logger {
	return ctx.logger.Named(name)
}

func (ctx *AppContext) SM() ServiceManager {
	return ctx.ServiceManager
}

func (ctx *AppContext) DM() DAOManager {
	return ctx.DAOManger
}

func (ctx *AppContext) CleanUp() {
	if err := ctx.logger.Sync(); err != nil {
		ctx.logger.Error("failed to sync logger", zap.Error(err))
	}
}
