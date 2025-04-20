package database

import (
	"github.com/carsonkrueger/main/interfaces"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

func GetByPK[T postgres.WritableTable, PK any, R any](DAO interfaces.IDatabaseObject[T, R], pk PK, db qrm.Queryable) error {
	return DAO.Table().
		SELECT(DAO.AllCols()).
		WHERE(DAO.PKColumn().EQ(pk)).
		Query(db, &R{})
}

func Insert[T postgres.WritableTable, R any](DAO interfaces.IDatabaseObject[T, R], values *R, db qrm.Queryable) error {
	return DAO.Table().
		INSERT(DAO.InsertCols()).
		VALUES(values).
		RETURNING(DAO.AllCols()).
		Query(db, values)
}

func InsertMany[T postgres.WritableTable, PK any, R any](DAO interfaces.IDatabaseObject[T, R], values []*R, db qrm.Queryable) error {
	return DAO.Table().
		INSERT(DAO.InsertCols()).
		MODELS(values).
		RETURNING(DAO.AllCols()).
		Query(db, values)
}

func Upsert[T postgres.WritableTable, R any](DAO interfaces.IDatabaseObject[T, R], values *R, db qrm.Queryable) error {
	return DAO.Table().
		INSERT(DAO.InsertCols()).
		VALUES(values).
		ON_CONFLICT(DAO.ConflictCols()...).
		DO_UPDATE(postgres.SET(DAO.UpdateOnConflictCols()...)).
		RETURNING(DAO.AllCols()).
		Query(db, values)
}

func UpsertMany[T postgres.WritableTable, R any](DAO interfaces.IDatabaseObject[T, R], values []*R, db qrm.Queryable) error {
	return DAO.Table().
		INSERT(DAO.InsertCols()).
		MODELS(values).
		ON_CONFLICT(DAO.ConflictCols()...).
		DO_UPDATE(postgres.SET(DAO.UpdateOnConflictCols()...)).
		RETURNING(DAO.AllCols()).
		Query(db, values)
}
