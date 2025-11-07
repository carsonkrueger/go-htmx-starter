package context

import (
	_ "github.com/lib/pq"
)

type AppContext struct {
	ServiceManager ServiceManager
	DAOManger      DAOManager
}

func NewAppContext(
	ServiceManager ServiceManager,
	DAOManger DAOManager,
) *AppContext {
	return &AppContext{
		ServiceManager,
		DAOManger,
	}
}

func (ctx *AppContext) SM() ServiceManager {
	return ctx.ServiceManager
}

func (ctx *AppContext) DM() DAOManager {
	return ctx.DAOManger
}
