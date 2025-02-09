package types

import "github.com/carsonkrueger/main/templates"

type AppContext struct {
	Templates *templates.Templates
}

type SetCtx interface {
	SetCtx(ctx *AppContext)
}

type WithContext struct {
	ctx *AppContext
}

func (b *WithContext) SetCtx(ctx *AppContext) {
	b.ctx = ctx
}

func (b *WithContext) GetCtx() *AppContext {
	return b.ctx
}
