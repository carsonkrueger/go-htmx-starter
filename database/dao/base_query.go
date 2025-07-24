package dao

import (
	gctx "context"
	"time"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/models"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

func index[PK context.PrimaryKey, R any](ctx gctx.Context, DAO context.DAO[PK, R], params *models.SearchParams, db qrm.Queryable) ([]*R, error) {
	query := DAO.Table().SELECT(DAO.AllCols())
	if params != nil {
		if params.Where != nil {
			query = query.WHERE(params.Where)
		}
		if len(params.OrderBy) > 0 {
			query = query.ORDER_BY(params.OrderBy...)
		}
		if params.Offset != nil {
			query = query.OFFSET(*params.Offset)
		}
		if params.Limit != nil {
			query = query.LIMIT(*params.Limit)
		}
	}
	models := []*R{}
	if err := query.QueryContext(ctx, db, &models); err != nil {
		return nil, err
	}
	return models, nil
}

func getOne[PK context.PrimaryKey, R any](ctx gctx.Context, DAO context.DAO[PK, R], pk PK, db qrm.Queryable) (*R, error) {
	var model R
	if err := DAO.Table().
		SELECT(DAO.AllCols()).
		WHERE(DAO.PKMatch(pk)).
		LIMIT(1).
		QueryContext(ctx, db, &model); err != nil {
		return nil, err
	}
	return &model, nil
}

func getMany[PK context.PrimaryKey, R any](ctx gctx.Context, DAO context.DAO[PK, R], pks []PK, db qrm.Queryable) ([]*R, error) {
	models := []*R{}
	if len(pks) == 0 {
		return models, nil
	}
	where := DAO.PKMatch(pks[0])
	for _, pk := range pks[1:] {
		where = where.OR(DAO.PKMatch(pk))
	}
	if err := DAO.Table().
		SELECT(DAO.AllCols()).
		WHERE(where).
		QueryContext(ctx, db, &models); err != nil {
		return nil, err
	}
	return models, nil
}

func insert[PK any, R any](ctx gctx.Context, DAO context.DAO[PK, R], model *R, db qrm.Queryable) error {
	return DAO.Table().
		INSERT(DAO.InsertCols()).
		MODEL(model).
		RETURNING(DAO.AllCols()).
		QueryContext(ctx, db, model)
}

func insertMany[PK any, R any](ctx gctx.Context, DAO context.DAO[PK, R], models *[]*R, db qrm.Queryable) error {
	return DAO.Table().
		INSERT(DAO.InsertCols()).
		MODELS(models).
		RETURNING(DAO.AllCols()).
		QueryContext(ctx, db, models)
}

func upsert[PK any, R any](ctx gctx.Context, DAO context.DAO[PK, R], model *R, db qrm.Queryable) error {
	up := DAO.GetUpdatedAt(model)
	if up != nil {
		*up = time.Now()
	}
	conflictCols := DAO.OnConflictCols()
	updateCols := DAO.UpdateOnConflictCols()
	query := DAO.Table().
		INSERT(DAO.UpdateCols()).
		MODEL(model)
	if len(updateCols) > 0 && len(conflictCols) > 0 {
		query = query.
			ON_CONFLICT(conflictCols...).
			DO_UPDATE(postgres.SET(updateCols...))
	}
	return query.
		RETURNING(DAO.AllCols()).
		QueryContext(ctx, db, model)
}

func upsertMany[PK any, R any](ctx gctx.Context, DAO context.DAO[PK, R], models *[]*R, db qrm.Queryable) error {
	for _, v := range *models {
		up := DAO.GetUpdatedAt(v)
		if up != nil {
			*up = time.Now()
		}
	}
	conflictCols := DAO.OnConflictCols()
	updateCols := DAO.UpdateOnConflictCols()
	query := DAO.Table().
		INSERT(DAO.UpdateCols()).
		MODELS(models)
	if len(updateCols) > 0 && len(conflictCols) > 0 {
		query = query.
			ON_CONFLICT(conflictCols...).
			DO_UPDATE(postgres.SET(updateCols...))
	}
	return query.
		RETURNING(DAO.AllCols()).
		QueryContext(ctx, db, models)
}

func update[PK any, R any](ctx gctx.Context, DAO context.DAO[PK, R], model *R, pk PK, db qrm.Queryable) error {
	up := DAO.GetUpdatedAt(model)
	if up != nil {
		*up = time.Now()
	}
	return DAO.Table().
		UPDATE(DAO.UpdateCols()).
		MODEL(model).
		WHERE(DAO.PKMatch(pk)).
		RETURNING(DAO.AllCols()).
		QueryContext(ctx, db, model)
}

func delete[PK any, R any](ctx gctx.Context, DAO context.DAO[PK, R], pk PK, db qrm.Executable) error {
	_, err := DAO.Table().
		DELETE().
		WHERE(DAO.PKMatch(pk)).
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

func (q *baseDAOQueryable[PK, R]) Index(ctx gctx.Context, params *models.SearchParams) ([]*R, error) {
	return index(ctx, q.Dao, params, context.GetDB(ctx))
}

func (q *baseDAOQueryable[PK, R]) GetOne(ctx gctx.Context, pk PK) (*R, error) {
	return getOne(ctx, q.Dao, pk, context.GetDB(ctx))
}

func (q *baseDAOQueryable[PK, R]) GetMany(ctx gctx.Context, pks []PK) ([]*R, error) {
	return getMany(ctx, q.Dao, pks, context.GetDB(ctx))
}

func (q *baseDAOQueryable[PK, R]) Insert(ctx gctx.Context, model *R) error {
	return insert(ctx, q.Dao, model, context.GetDB(ctx))
}

func (q *baseDAOQueryable[PK, R]) InsertMany(ctx gctx.Context, models *[]*R) error {
	return insertMany(ctx, q.Dao, models, context.GetDB(ctx))
}

func (q *baseDAOQueryable[PK, R]) Upsert(ctx gctx.Context, model *R) error {
	return upsert(ctx, q.Dao, model, context.GetDB(ctx))
}

func (q *baseDAOQueryable[PK, R]) UpsertMany(ctx gctx.Context, models *[]*R) error {
	return upsertMany(ctx, q.Dao, models, context.GetDB(ctx))
}

func (q *baseDAOQueryable[PK, R]) Update(ctx gctx.Context, model *R, pk PK) error {
	return update(ctx, q.Dao, model, pk, context.GetDB(ctx))
}

func (q *baseDAOQueryable[PK, R]) Delete(ctx gctx.Context, pk PK) error {
	return delete(ctx, q.Dao, pk, context.GetDB(ctx))
}
