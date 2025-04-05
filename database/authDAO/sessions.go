package authDAO

import (
	"database/sql"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/go-jet/jet/v2/postgres"
)

type sessionsDAO struct {
	db *sql.DB
}

func NewSessionsDAO(db *sql.DB) *sessionsDAO {
	return &sessionsDAO{
		db,
	}
}

func (dao *sessionsDAO) GetById(id authModels.SessionsPrimaryKey) (*model.Sessions, error) {
	var row model.Sessions
	err := table.Sessions.
		SELECT(table.Sessions.AllColumns).
		FROM(table.Sessions).
		WHERE(table.Sessions.UserID.EQ(postgres.Int(id.UserID)).
			AND(table.Sessions.Token.EQ(postgres.String(id.AuthToken)))).
		LIMIT(1).
		Query(dao.db, &row)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (dao *sessionsDAO) Insert(row *model.Sessions) error {
	return table.Sessions.
		INSERT(table.Sessions.AllColumns.Except(table.Sessions.CreatedAt)).
		VALUES(row.UserID, row.Token, row.ExpiresAt).
		RETURNING(table.Sessions.UserID, table.Sessions.Token).
		Query(dao.db, row)
}

func (dao *sessionsDAO) InsertMany(rows []*model.Sessions) error {
	if len(rows) == 0 {
		return nil
	}
	return table.Sessions.
		INSERT(table.Sessions.AllColumns.Except(table.Sessions.CreatedAt)).
		MODELS(rows).
		RETURNING(table.Sessions.UserID, table.Sessions.Token).
		Query(dao.db, &rows)
}

// Returns ID int64 if inserted.
// Parameter cols_update are the columns to be updated on conflict - if not provided, a few columns are updated
func (dao *sessionsDAO) Upsert(row *model.Sessions, colsUpdate ...postgres.ColumnAssigment) error {
	if len(colsUpdate) == 0 {
		colsUpdate = []postgres.ColumnAssigment{
			table.Sessions.Token.SET(table.Sessions.Token),
			table.Sessions.ExpiresAt.SET(table.Sessions.ExpiresAt),
		}
	}

	return table.Sessions.
		INSERT(table.Sessions.AllColumns.Except(table.Sessions.CreatedAt)).
		VALUES(row.UserID, row.Token, row.ExpiresAt).
		ON_CONFLICT(table.Sessions.Token).
		DO_UPDATE(postgres.SET(colsUpdate...)).
		RETURNING(table.Sessions.UserID, table.Sessions.Token).
		Query(dao.db, row)
}

func (dao *sessionsDAO) UpsertMany(rows []*model.Sessions, colsUpdate ...postgres.ColumnAssigment) error {
	if len(colsUpdate) == 0 {
		colsUpdate = []postgres.ColumnAssigment{
			table.Sessions.Token.SET(table.Sessions.Token),
			table.Sessions.ExpiresAt.SET(table.Sessions.ExpiresAt),
		}
	}

	return table.Sessions.
		INSERT(table.Sessions.AllColumns.Except(table.Sessions.CreatedAt)).
		MODELS(rows).
		ON_CONFLICT(table.Sessions.Token).
		DO_UPDATE(postgres.SET(colsUpdate...)).
		RETURNING(table.Sessions.UserID, table.Sessions.Token).
		Query(dao.db, rows)
}

func (dao *sessionsDAO) Update(row *model.Sessions) error {
	_, err := table.Sessions.
		UPDATE(table.Sessions.AllColumns.Except(table.Sessions.CreatedAt)).
		MODEL(row).
		WHERE(table.Sessions.UserID.EQ(postgres.Int(row.UserID)).
			AND(table.Sessions.Token.EQ(postgres.String(row.Token)))).
		Exec(dao.db)
	return err
}

func (dao *sessionsDAO) Delete(id authModels.SessionsPrimaryKey) error {
	_, err := table.Sessions.DELETE().
		WHERE(table.Sessions.UserID.EQ(postgres.Int(id.UserID)).
			AND(table.Sessions.Token.EQ(postgres.String(id.AuthToken)))).
		Exec(dao.db)
	if err != nil {
		return err
	}
	return nil
}

func (dao *sessionsDAO) GetAll() (*[]model.Sessions, error) {
	var rows []model.Sessions
	err := table.Sessions.
		SELECT(table.Sessions.AllColumns).
		ORDER_BY(table.Sessions.UserID.DESC()).
		Query(dao.db, &rows)
	if err != nil {
		return nil, err
	}
	return &rows, nil
}
