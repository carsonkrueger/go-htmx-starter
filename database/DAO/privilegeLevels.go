package DAO

import (
	"database/sql"
	"time"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/go-jet/jet/v2/postgres"
)

type privilegeLevelsDAO struct {
	db *sql.DB
	DAOBaseQueries[int64, model.PrivilegeLevels]
}

func newPrivilegeLevelsDAO(db *sql.DB) *privilegeLevelsDAO {
	dao := &privilegeLevelsDAO{
		db:             db,
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[int64, model.PrivilegeLevels](dao)
	dao.DAOBaseQueries = &queries
	return dao
}

func (dao *privilegeLevelsDAO) Table() PostgresTable {
	return table.PrivilegeLevels
}

func (dao *privilegeLevelsDAO) InsertCols() postgres.ColumnList {
	return table.PrivilegeLevels.AllColumns.Except(
		table.PrivilegeLevels.ID,
		table.PrivilegeLevels.CreatedAt,
		table.PrivilegeLevels.UpdatedAt,
	)
}

func (dao *privilegeLevelsDAO) UpdateCols() postgres.ColumnList {
	return table.PrivilegeLevels.AllColumns.Except(
		table.PrivilegeLevels.ID,
		table.PrivilegeLevels.CreatedAt,
	)
}

func (dao *privilegeLevelsDAO) AllCols() postgres.ColumnList {
	return table.PrivilegeLevels.AllColumns
}

func (dao *privilegeLevelsDAO) OnConflictCols() postgres.ColumnList {
	return []postgres.Column{table.PrivilegeLevels.Name}
}

func (dao *privilegeLevelsDAO) UpdateOnConflictCols() []postgres.ColumnAssigment {
	return []postgres.ColumnAssigment{
		table.PrivilegeLevels.Name.SET(table.PrivilegeLevels.EXCLUDED.Name),
	}
}

func (dao *privilegeLevelsDAO) PKMatch(pk int64) postgres.BoolExpression {
	return table.PrivilegeLevels.ID.EQ(postgres.Int(pk))
}

func (dao *privilegeLevelsDAO) GetUpdatedAt(row *model.PrivilegeLevels) *time.Time {
	return row.UpdatedAt
}
