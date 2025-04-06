package context

import (
	gctx "context"

	"github.com/carsonkrueger/main/constant"
)

func WithToken(ctx gctx.Context, token string) gctx.Context {
	return gctx.WithValue(ctx, constant.AUTH_TOKEN_KEY, token)
}

func GetToken(ctx gctx.Context) string {
	return ctx.Value(constant.AUTH_TOKEN_KEY).(string)
}

var USER_ID_KEY = "USER_ID"

func WithUserId(ctx gctx.Context, id int64) gctx.Context {
	return gctx.WithValue(ctx, USER_ID_KEY, id)
}

func GetUserId(ctx gctx.Context) int64 {
	return ctx.Value(USER_ID_KEY).(int64)
}

var PRIVILEGE_LEVEL_ID_KEY = "PRIVILEGE_LEVEL_ID"

func WithPrivilegeLevelID(ctx gctx.Context, id int64) gctx.Context {
	return gctx.WithValue(ctx, PRIVILEGE_LEVEL_ID_KEY, id)
}

func GetPrivilegeLevelID(ctx gctx.Context) int64 {
	return ctx.Value(PRIVILEGE_LEVEL_ID_KEY).(int64)
}
