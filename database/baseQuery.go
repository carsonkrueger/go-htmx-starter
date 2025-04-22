package database

import (
	"time"

	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/models"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

func index[PK interfaces.PrimaryKey, R any](DAO interfaces.IDAO[PK, R], params *models.SearchParams, db qrm.Queryable) ([]*R, error) {
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
	var models []*R
	if err := query.Query(db, &models); err != nil {
		return nil, err
	}
	return models, nil
}

func getOne[PK interfaces.PrimaryKey, R any](DAO interfaces.IDAO[PK, R], pk PK, db qrm.Queryable) (*R, error) {
	var model R
	if err := DAO.Table().
		SELECT(DAO.AllCols()).
		WHERE(DAO.PKMatch(pk)).
		LIMIT(1).
		Query(db, &model); err != nil {
		return nil, err
	}
	return &model, nil
}

func getMany[PK interfaces.PrimaryKey, R any](DAO interfaces.IDAO[PK, R], pk PK, db qrm.Queryable) ([]*R, error) {
	var models []*R
	if err := DAO.Table().
		SELECT(DAO.AllCols()).
		WHERE(DAO.PKMatch(pk)).
		Query(db, &models); err != nil {
		return nil, err
	}
	return models, nil
}

func insert[PK any, R any](DAO interfaces.IDAO[PK, R], model *R, db qrm.Queryable) error {
	return DAO.Table().
		INSERT(DAO.InsertCols()).
		MODEL(model).
		RETURNING(DAO.AllCols()).
		Query(db, model)
}

func insertMany[PK any, R any](DAO interfaces.IDAO[PK, R], models *[]*R, db qrm.Queryable) error {
	return DAO.Table().
		INSERT(DAO.InsertCols()).
		MODELS(models).
		RETURNING(DAO.AllCols()).
		Query(db, &models)
}

func upsert[PK any, R any](DAO interfaces.IDAO[PK, R], model *R, db qrm.Queryable) error {
	up := DAO.GetUpdatedAt(model)
	if up != nil {
		*up = time.Now()
	}
	return DAO.Table().
		INSERT(DAO.UpdateCols()).
		MODEL(model).
		ON_CONFLICT(DAO.OnConflictCols()...).
		DO_UPDATE(postgres.SET(DAO.UpdateOnConflictCols()...)).
		RETURNING(DAO.AllCols()).
		Query(db, model)
}

func upsertMany[PK any, R any](DAO interfaces.IDAO[PK, R], models *[]*R, db qrm.Queryable) error {
	for _, v := range *models {
		up := DAO.GetUpdatedAt(v)
		if up != nil {
			*up = time.Now()
		}
	}
	return DAO.Table().
		INSERT(DAO.UpdateCols()).
		MODELS(*models).
		ON_CONFLICT(DAO.OnConflictCols()...).
		DO_UPDATE(postgres.SET(DAO.UpdateOnConflictCols()...)).
		RETURNING(DAO.AllCols()).
		Query(db, models)
}

func update[PK any, R any](DAO interfaces.IDAO[PK, R], model *R, pk PK, db qrm.Queryable) error {
	up := DAO.GetUpdatedAt(model)
	if up != nil {
		*up = time.Now()
	}
	return DAO.Table().
		UPDATE(DAO.UpdateCols()).
		MODEL(model).
		WHERE(DAO.PKMatch(pk)).
		RETURNING(DAO.AllCols()).
		Query(db, model)
}

func delete[PK any, R any](DAO interfaces.IDAO[PK, R], pk PK, db qrm.Executable) error {
	_, err := DAO.Table().
		DELETE().
		WHERE(DAO.PKMatch(pk)).
		Exec(db)
	return err
}

type baseDAOQueryable[PK interfaces.PrimaryKey, R any] struct {
	Dao interfaces.IDAO[PK, R]
}

func NewDAOQueryable[PK interfaces.PrimaryKey, R any](dao interfaces.IDAO[PK, R]) baseDAOQueryable[PK, R] {
	return baseDAOQueryable[PK, R]{
		dao,
	}
}

func (q *baseDAOQueryable[PK, R]) Index(params *models.SearchParams, db qrm.Queryable) ([]*R, error) {
	return index(q.Dao, params, db)
}

func (q *baseDAOQueryable[PK, R]) GetOne(pk PK, db qrm.Queryable) (*R, error) {
	return getOne(q.Dao, pk, db)
}

func (q *baseDAOQueryable[PK, R]) GetMany(pk PK, db qrm.Queryable) ([]*R, error) {
	return getMany(q.Dao, pk, db)
}

func (q *baseDAOQueryable[PK, R]) Insert(model *R, db qrm.Queryable) error {
	return insert(q.Dao, model, db)
}

func (q *baseDAOQueryable[PK, R]) InsertMany(models *[]*R, db qrm.Queryable) error {
	return insertMany(q.Dao, models, db)
}

func (q *baseDAOQueryable[PK, R]) Upsert(model *R, db qrm.Queryable) error {
	return upsert(q.Dao, model, db)
}

func (q *baseDAOQueryable[PK, R]) UpsertMany(models *[]*R, db qrm.Queryable) error {
	return upsertMany(q.Dao, models, db)
}

func (q *baseDAOQueryable[PK, R]) Update(model *R, pk PK, db qrm.Queryable) error {
	return update(q.Dao, model, pk, db)
}

func (q *baseDAOQueryable[PK, R]) Delete(pk PK, db qrm.Executable) error {
	return delete(q.Dao, pk, db)
}
