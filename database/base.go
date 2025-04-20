package database

import (
	"time"

	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/models"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

func Index[PK interfaces.PK, R any](DAO interfaces.IDAO[PK, R], params models.SearchParams, db qrm.Queryable) ([]*R, error) {
	var rows []*R
	query := DAO.Table().SELECT(DAO.AllCols())
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
	if err := query.Query(db, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func GetOne[PK interfaces.PK, R any](DAO interfaces.IDAO[PK, R], pk PK, model *R, db qrm.Queryable) error {
	return DAO.Table().
		SELECT(DAO.AllCols()).
		WHERE(DAO.PKMatch(pk)).
		LIMIT(1).
		Query(db, model)
}

func GetMany[PK interfaces.PK, R any](DAO interfaces.IDAO[PK, R], pk PK, models *[]*R, db qrm.Queryable) error {
	return DAO.Table().
		SELECT(DAO.AllCols()).
		WHERE(DAO.PKMatch(pk)).
		Query(db, &models)
}

func Insert[PK any, R any](DAO interfaces.IDAO[PK, R], model *R, db qrm.Queryable) error {
	return DAO.Table().
		INSERT(DAO.InsertCols()).
		MODEL(model).
		RETURNING(DAO.AllCols()).
		Query(db, model)
}

func InsertMany[T interfaces.IPostgresTable, PK any, R any](DAO interfaces.IDAO[PK, R], models []*R, db qrm.Queryable) error {
	return DAO.Table().
		INSERT(DAO.InsertCols()).
		MODELS(models).
		RETURNING(DAO.AllCols()).
		Query(db, &models)
}

func Upsert[PK any, R any](DAO interfaces.IDAO[PK, R], model *R, db qrm.Queryable) error {
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

func UpsertMany[PK any, R any](DAO interfaces.IDAO[PK, R], models []*R, db qrm.Queryable) error {
	for _, v := range models {
		up := DAO.GetUpdatedAt(v)
		if up != nil {
			*up = time.Now()
		}
	}
	return DAO.Table().
		INSERT(DAO.UpdateCols()).
		MODELS(models).
		ON_CONFLICT(DAO.OnConflictCols()...).
		DO_UPDATE(postgres.SET(DAO.UpdateOnConflictCols()...)).
		RETURNING(DAO.AllCols()).
		Query(db, &models)
}

func Update[PK any, R any](DAO interfaces.IDAO[PK, R], model *R, pk PK, db qrm.Queryable) error {
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

func Delete[PK any, R any](DAO interfaces.IDAO[PK, R], pk PK, db qrm.Executable) error {
	_, err := DAO.Table().
		DELETE().
		WHERE(DAO.PKMatch(pk)).
		Exec(db)
	return err
}
