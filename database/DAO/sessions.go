package DAO

import (
	"database/sql"
	"time"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/go-jet/jet/v2/postgres"
)

type sessionsDAO struct {
	db *sql.DB
	DAOBaseQueries[authModels.SessionsPrimaryKey, model.Sessions]
}

func newSessionsDAO(db *sql.DB) *sessionsDAO {
	dao := &sessionsDAO{
		db:             db,
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[authModels.SessionsPrimaryKey, model.Sessions](dao)
	dao.DAOBaseQueries = &queries
	return dao
}

func (dao *sessionsDAO) Table() PostgresTable {
	return table.Sessions
}

func (dao *sessionsDAO) InsertCols() postgres.ColumnList {
	return table.Sessions.AllColumns.Except(
		table.Sessions.CreatedAt,
	)
}

func (dao *sessionsDAO) UpdateCols() postgres.ColumnList {
	return table.Sessions.AllColumns.Except(
		table.Sessions.CreatedAt,
		table.Sessions.UserID,
		table.Sessions.Token,
	)
}

func (dao *sessionsDAO) AllCols() postgres.ColumnList {
	return table.Sessions.AllColumns
}

func (dao *sessionsDAO) OnConflictCols() postgres.ColumnList {
	return []postgres.Column{}
}

func (dao *sessionsDAO) UpdateOnConflictCols() []postgres.ColumnAssigment {
	return []postgres.ColumnAssigment{}
}

func (dao *sessionsDAO) PKMatch(pk authModels.SessionsPrimaryKey) postgres.BoolExpression {
	return table.Sessions.
		UserID.EQ(postgres.Int(pk.UserID)).
		AND(table.Sessions.Token.EQ(postgres.String(pk.AuthToken)))
}

func (dao *sessionsDAO) GetUpdatedAt(row *model.Sessions) *time.Time {
	return nil
}
