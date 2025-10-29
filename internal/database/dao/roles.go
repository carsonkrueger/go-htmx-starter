package dao

import (
	"time"

	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/database/gen/go_starter_db/auth/table"
	"github.com/carsonkrueger/main/pkg/db/auth/model"
	dbmodel "github.com/carsonkrueger/main/pkg/db/auth/model"
	"github.com/go-jet/jet/v2/postgres"
)

type RolesDAO struct {
	context.DAOBaseQueries[int16, model.Roles]
}

func NewRolesDAO() *RolesDAO {
	dao := &RolesDAO{
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[int16, dbmodel.Roles](dao)
	dao.DAOBaseQueries = &queries
	return dao
}

func (dao *RolesDAO) Table() context.PostgresTable {
	return table.Roles
}

func (dao *RolesDAO) InsertCols() postgres.ColumnList {
	return table.Roles.AllColumns.Except(
		table.Roles.ID,
		table.Roles.CreatedAt,
		table.Roles.UpdatedAt,
	)
}

func (dao *RolesDAO) UpdateCols() postgres.ColumnList {
	return table.Roles.AllColumns.Except(
		table.Roles.ID,
		table.Roles.CreatedAt,
	)
}

func (dao *RolesDAO) AllCols() postgres.ColumnList {
	return table.Roles.AllColumns
}

func (dao *RolesDAO) OnConflictCols() postgres.ColumnList {
	return []postgres.Column{table.Roles.Name}
}

func (dao *RolesDAO) UpdateOnConflictCols() []postgres.ColumnAssigment {
	return []postgres.ColumnAssigment{
		table.Roles.Name.SET(table.Roles.EXCLUDED.Name),
	}
}

func (dao *RolesDAO) PKMatch(pk int16) postgres.BoolExpression {
	return table.Roles.ID.EQ(postgres.Int16(pk))
}

func (dao *RolesDAO) GetUpdatedAt(row *dbmodel.Roles) *time.Time {
	return row.UpdatedAt
}
