package context

import (
	gctx "context"
	"database/sql"
	"errors"

	"github.com/carsonkrueger/main/constant"
	"github.com/go-jet/jet/v2/qrm"
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

var DB_CONNECTION_KEY = "DB_CONNECTION"

func WithDB(ctx gctx.Context, db *sql.DB) gctx.Context {
	return gctx.WithValue(ctx, DB_CONNECTION_KEY, db)
}

func GetDB(ctx gctx.Context) qrm.DB {
	db := ctx.Value(DB_CONNECTION_KEY)
	switch sdb := db.(type) {
	case qrm.DB:
		return sdb
	default:
		return nil
	}
}

func BeginTx(ctx gctx.Context) (*sql.Tx, error) {
	switch db := GetDB(ctx).(type) {
	case *sql.DB:
		return db.BeginTx(ctx, nil)
	case *sql.Tx:
		return db, nil
	default:
		return nil, errors.New("unsupported database type")
	}
}
