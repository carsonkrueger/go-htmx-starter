package types

type AppContext struct{}

type SetCtx interface {
	SetCtx(ctx *AppContext)
}

type WithContext struct {
	Ctx *AppContext
}

func (b *WithContext) SetCtx(ctx *AppContext) {
	b.Ctx = ctx
}
