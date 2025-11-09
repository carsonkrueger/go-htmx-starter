package dao

import (
	gctx "context"
	"time"

	"github.com/carsonkrueger/main/internal/context"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

func getAll[PK context.PrimaryKey, R any, D context.DAO[PK, R]](ctx gctx.Context, dao D, db qrm.Queryable) ([]R, error) {
	query := dao.Table().SELECT(dao.AllCols())
	models := []R{}
	if err := query.QueryContext(ctx, db, &models); err != nil {
		return nil, err
	}
	return models, nil
}

func getOne[PK context.PrimaryKey, R any, D context.DAO[PK, R]](ctx gctx.Context, dao D, pk PK, db qrm.Queryable) (R, error) {
	var model R
	if err := dao.Table().
		SELECT(dao.AllCols()).
		WHERE(dao.PKMatch(pk)).
		LIMIT(1).
		QueryContext(ctx, db, &model); err != nil {
		return model, err
	}
	return model, nil
}

func getMany[PK context.PrimaryKey, R any, D context.DAO[PK, R]](ctx gctx.Context, DAO D, pks []PK, db qrm.Queryable) ([]R, error) {
	models := []R{}
	if len(pks) == 0 {
		return models, nil
	}
	where := DAO.PKMatch(pks[0])
	for _, pk := range pks[1:] {
		where = where.OR(DAO.PKMatch(pk))
	}
	err := DAO.Table().
		SELECT(DAO.AllCols()).
		WHERE(where).
		QueryContext(ctx, db, &models)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func insert[PK any, R any, D context.DAO[PK, R]](ctx gctx.Context, dao D, model *R, db qrm.Queryable) error {
	return dao.Table().
		INSERT(dao.InsertCols()).
		MODEL(model).
		RETURNING(dao.AllCols()).
		QueryContext(ctx, db, model)
}

func insertMany[PK any, R any, D context.DAO[PK, R]](ctx gctx.Context, dao D, models []R, db qrm.Queryable) error {
	return dao.Table().
		INSERT(dao.InsertCols()).
		MODELS(models).
		RETURNING(dao.AllCols()).
		QueryContext(ctx, db, &models)
}

func upsert[PK any, R any, D context.DAO[PK, R]](ctx gctx.Context, dao D, model *R, db qrm.Queryable) error {
	up := dao.GetUpdatedAt(model)
	if up != nil {
		*up = time.Now()
	}
	conflictCols := dao.OnConflictCols()
	updateCols := dao.UpdateOnConflictCols()
	query := dao.Table().
		INSERT(dao.UpdateCols()).
		MODEL(model)
	if len(updateCols) > 0 && len(conflictCols) > 0 {
		query = query.
			ON_CONFLICT(conflictCols...).
			DO_UPDATE(postgres.SET(updateCols...))
	}
	return query.
		RETURNING(dao.AllCols()).
		QueryContext(ctx, db, model)
}

func upsertMany[PK any, R any, D context.DAO[PK, R]](ctx gctx.Context, dao D, models []R, db qrm.Queryable) error {
	if len(models) == 0 {
		return nil
	}
	for _, v := range models {
		up := dao.GetUpdatedAt(&v)
		if up != nil {
			*up = time.Now()
		}
	}
	conflictCols := dao.OnConflictCols()
	updateCols := dao.UpdateOnConflictCols()
	query := dao.Table().
		INSERT(dao.UpdateCols()).
		MODELS(&models)
	if len(updateCols) > 0 && len(conflictCols) > 0 {
		query = query.
			ON_CONFLICT(conflictCols...).
			DO_UPDATE(postgres.SET(updateCols...))
	}
	return query.
		RETURNING(dao.AllCols()).
		QueryContext(ctx, db, &models)
}

func update[PK any, R any, D context.DAO[PK, R]](ctx gctx.Context, dao D, model *R, pk PK, db qrm.Queryable) error {
	up := dao.GetUpdatedAt(model)
	if up != nil {
		*up = time.Now()
	}
	return dao.Table().
		UPDATE(dao.UpdateCols()).
		MODEL(model).
		WHERE(dao.PKMatch(pk)).
		RETURNING(dao.AllCols()).
		QueryContext(ctx, db, model)
}

func delete[PK any, R any, D context.DAO[PK, R]](ctx gctx.Context, dao D, pk PK, db qrm.Executable) error {
	_, err := dao.Table().
		DELETE().
		WHERE(dao.PKMatch(pk)).
		ExecContext(ctx, db)
	return err
}

type baseDAOQueryable[PK context.PrimaryKey, R any] struct {
	Dao context.DAO[PK, R]
}

func newDAOQueryable[PK context.PrimaryKey, R any](dao context.DAO[PK, R]) baseDAOQueryable[PK, R] {
	return baseDAOQueryable[PK, R]{
		dao,
	}
}

func (q *baseDAOQueryable[PK, R]) GetAll(ctx gctx.Context) ([]R, error) {
	return getAll(ctx, q.Dao, context.GetDB(ctx))
}

func (q *baseDAOQueryable[PK, R]) GetOne(ctx gctx.Context, pk PK) (R, error) {
	return getOne(ctx, q.Dao, pk, context.GetDB(ctx))
}

func (q *baseDAOQueryable[PK, R]) GetMany(ctx gctx.Context, pks []PK) ([]R, error) {
	return getMany(ctx, q.Dao, pks, context.GetDB(ctx))
}

func (q *baseDAOQueryable[PK, R]) Insert(ctx gctx.Context, model *R) error {
	return insert(ctx, q.Dao, model, context.GetDB(ctx))
}

func (q *baseDAOQueryable[PK, R]) InsertMany(ctx gctx.Context, models []R) error {
	return insertMany(ctx, q.Dao, models, context.GetDB(ctx))
}

func (q *baseDAOQueryable[PK, R]) Upsert(ctx gctx.Context, model *R) error {
	return upsert(ctx, q.Dao, model, context.GetDB(ctx))
}

func (q *baseDAOQueryable[PK, R]) UpsertMany(ctx gctx.Context, models []R) error {
	return upsertMany(ctx, q.Dao, models, context.GetDB(ctx))
}

func (q *baseDAOQueryable[PK, R]) Update(ctx gctx.Context, model *R, pk PK) error {
	return update(ctx, q.Dao, model, pk, context.GetDB(ctx))
}

func (q *baseDAOQueryable[PK, R]) Delete(ctx gctx.Context, pk PK) error {
	return delete(ctx, q.Dao, pk, context.GetDB(ctx))
}
