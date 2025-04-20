package authDAO

import (
	"database/sql"
	"time"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/carsonkrueger/main/interfaces"
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

func (dao *sessionsDAO) Table() interfaces.IPostgresTable {
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
