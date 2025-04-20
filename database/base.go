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

func GetOne[PK interfaces.PK, R any](DAO interfaces.IDAO[PK, R], pk PK, db qrm.Queryable) (*R, error) {
	var row *R
	err := DAO.Table().
		SELECT(DAO.AllCols()).
		WHERE(DAO.PKMatch(pk)).
		LIMIT(1).
		Query(db, row)
	if row == nil {
		return nil, err
	}
	return row, nil
}

func GetMany[PK interfaces.PK, R any](DAO interfaces.IDAO[PK, R], pk PK, db qrm.Queryable) ([]*R, error) {
	var rows []*R
	err := DAO.Table().
		SELECT(DAO.AllCols()).
		WHERE(DAO.PKMatch(pk)).
		Query(db, &rows)
	if len(rows) == 0 {
		return nil, err
	}
	return rows, nil
}

func Insert[PK any, R any](DAO interfaces.IDAO[PK, R], values *R, db qrm.Queryable) error {
	return DAO.Table().
		INSERT(DAO.InsertCols()).
		VALUES(values).
		RETURNING(DAO.AllCols()).
		Query(db, values)
}

func InsertMany[T interfaces.IPostgresTable, PK any, R any](DAO interfaces.IDAO[PK, R], values []*R, db qrm.Queryable) error {
	return DAO.Table().
		INSERT(DAO.InsertCols()).
		MODELS(values).
		RETURNING(DAO.AllCols()).
		Query(db, values)
}

func Upsert[PK any, R any](DAO interfaces.IDAO[PK, R], values *R, db qrm.Queryable) error {
	*DAO.GetUpdatedAt(values) = time.Now()
	return DAO.Table().
		INSERT(DAO.UpdateCols()).
		VALUES(values).
		ON_CONFLICT(DAO.OnConflictCols()...).
		DO_UPDATE(postgres.SET(DAO.UpdateOnConflictCols()...)).
		RETURNING(DAO.AllCols()).
		Query(db, values)
}

func UpsertMany[PK any, R any](DAO interfaces.IDAO[PK, R], values []*R, db qrm.Queryable) error {
	for _, v := range values {
		*DAO.GetUpdatedAt(v) = time.Now()
	}
	return DAO.Table().
		INSERT(DAO.UpdateCols()).
		MODELS(values).
		ON_CONFLICT(DAO.OnConflictCols()...).
		DO_UPDATE(postgres.SET(DAO.UpdateOnConflictCols()...)).
		RETURNING(DAO.AllCols()).
		Query(db, values)
}

func Update[PK any, R any](DAO interfaces.IDAO[PK, R], values *R, pk PK, db qrm.Queryable) error {
	*DAO.GetUpdatedAt(values) = time.Now()
	return DAO.Table().
		UPDATE(DAO.UpdateCols()).
		MODEL(values).
		WHERE(DAO.PKMatch(pk)).
		RETURNING(DAO.AllCols()).
		Query(db, values)
}

func Delete[PK any, R any](DAO interfaces.IDAO[PK, R], pk PK, db qrm.Executable) error {
	_, err := DAO.Table().
		DELETE().
		WHERE(DAO.PKMatch(pk)).
		Exec(db)
	return err
}
