package database

import (
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/models"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

func Index[T interfaces.IPostgresTable, R any](DAO interfaces.IDatabaseObject[T, R], params models.SearchParams, db qrm.Queryable) ([]*R, error) {
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
	err := query.Query(db, &rows)
	if len(rows) == 0 {
		return nil, err
	}
	return rows, nil
}

func GetOne[T interfaces.IPostgresTable, PK interfaces.PK, R any](DAO interfaces.IDatabaseObject[T, R], pk PK, db qrm.Queryable) (*R, error) {
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

func GetMany[T interfaces.IPostgresTable, PK interfaces.PK, R any](DAO interfaces.IDatabaseObject[T, R], pk PK, db qrm.Queryable) ([]*R, error) {
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

func Insert[T interfaces.IPostgresTable, R any](DAO interfaces.IDatabaseObject[T, R], values *R, db qrm.Queryable) error {
	return DAO.Table().
		INSERT(DAO.InsertCols()).
		VALUES(values).
		RETURNING(DAO.AllCols()).
		Query(db, values)
}

func InsertMany[T interfaces.IPostgresTable, PK any, R any](DAO interfaces.IDatabaseObject[T, R], values []*R, db qrm.Queryable) error {
	return DAO.Table().
		INSERT(DAO.InsertCols()).
		MODELS(values).
		RETURNING(DAO.AllCols()).
		Query(db, values)
}

func Upsert[T interfaces.IPostgresTable, R any](DAO interfaces.IDatabaseObject[T, R], values *R, db qrm.Queryable) error {
	return DAO.Table().
		INSERT(DAO.InsertCols()).
		VALUES(values).
		ON_CONFLICT(DAO.ConflictCols()...).
		DO_UPDATE(postgres.SET(DAO.UpdateOnConflictCols()...)).
		RETURNING(DAO.AllCols()).
		Query(db, values)
}

func UpsertMany[T interfaces.IPostgresTable, R any](DAO interfaces.IDatabaseObject[T, R], values []*R, db qrm.Queryable) error {
	return DAO.Table().
		INSERT(DAO.InsertCols()).
		MODELS(values).
		ON_CONFLICT(DAO.ConflictCols()...).
		DO_UPDATE(postgres.SET(DAO.UpdateOnConflictCols()...)).
		RETURNING(DAO.AllCols()).
		Query(db, values)
}

func Update[T interfaces.IPostgresTable, R any](DAO interfaces.IDatabaseObject[T, R], values *R, db qrm.Queryable) error {
	return DAO.Table().
		UPDATE(DAO.InsertCols()).
		MODEL(values).
		RETURNING(DAO.AllCols()).
		Query(db, values)
}

func Delete[T interfaces.IPostgresTable, PK interfaces.PK](DAO interfaces.IDatabaseObject[T, PK], pk PK, db qrm.Executable) error {
	_, err := DAO.Table().
		DELETE().
		WHERE(DAO.PKMatch(pk)).
		Exec(db)
	return err
}
