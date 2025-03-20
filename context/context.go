package context

import (
	gctx "context"

	"github.com/carsonkrueger/main/tools"
)

func WithToken(ctx gctx.Context, token string) gctx.Context {
	return gctx.WithValue(ctx, tools.AUTH_TOKEN_KEY, token)
}

func GetToken(ctx gctx.Context) string {
	return ctx.Value(tools.AUTH_TOKEN_KEY).(string)
}

var USER_ID_KEY = "USER_ID"

func WithUserId(ctx gctx.Context, id int64) gctx.Context {
	return gctx.WithValue(ctx, USER_ID_KEY, id)
}

func GetUserId(ctx gctx.Context) int64 {
	return ctx.Value(USER_ID_KEY).(int64)
}

var PRIVILEGE_ID_KEY = "PRIVILEGE_ID"

func WithPrivilegeID(ctx gctx.Context, id int64) gctx.Context {
	return gctx.WithValue(ctx, PRIVILEGE_ID_KEY, id)
}

func GetPrivilegeID(ctx gctx.Context) int64 {
	return ctx.Value(PRIVILEGE_ID_KEY).(int64)
}
