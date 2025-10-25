package context

import (
	gctx "context"
	"database/sql"
	"errors"

	"github.com/carsonkrueger/main/internal/constant"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	ErrTransactionAlreadyStarted = errors.New("context transaction has already started")
	ErrUnsupportedDatabase       = errors.New("unsupported context database")
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
	id, _ := ctx.Value(USER_ID_KEY).(int64)
	return id
}

var ROLE_ID_KEY = "ROLE_ID"

func WithRoleID(ctx gctx.Context, id int16) gctx.Context {
	return gctx.WithValue(ctx, ROLE_ID_KEY, id)
}

func GetRoleID(ctx gctx.Context) int16 {
	id, _ := ctx.Value(ROLE_ID_KEY).(int16)
	return id
}

var DB_CONNECTION_KEY = "DB_CONNECTION"

func WithDB(ctx gctx.Context, db qrm.DB) gctx.Context {
	return gctx.WithValue(ctx, DB_CONNECTION_KEY, db)
}

func GetDB(ctx gctx.Context) qrm.DB {
	return ctx.Value(DB_CONNECTION_KEY).(qrm.DB)
}

// Returns a new context that contains the transaction. Caller must Rollback and Commit manually. The new tx ctx cannot and should NOT be used after a rollback or commit.
//
// Suggested usage:
// ctx, tx, err := BeginTx(ctx)
//
//	if err != nil {
//	    return err
//	}
//
// defer tx.Rollback()
//
// // Perform database operations using new ctx
//
// tx.Commit()
func BeginTx(ctx gctx.Context) (gctx.Context, *sql.Tx, error) {
	switch db := GetDB(ctx).(type) {
	case *sql.DB:
		tx, err := db.BeginTx(ctx, nil)
		if err == nil {
			ctx = WithDB(ctx, tx)
		}
		return ctx, tx, err
	case *sql.Tx:
		return ctx, db, ErrTransactionAlreadyStarted
	default:
		return ctx, nil, ErrUnsupportedDatabase
	}
}
