package dao

import (
	"time"

	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/database/gen/go_starter_db/auth/table"
	dbmodel "github.com/carsonkrueger/main/pkg/db/auth/model"
	"github.com/carsonkrueger/main/pkg/model"
	"github.com/go-jet/jet/v2/postgres"
)

type RolesPrivilegesDAO struct {
	context.DAOBaseQueries[model.RolesPrivilegesPrimaryKey, dbmodel.RolesPrivileges]
}

func NewRolesPrivilegesDAO() *RolesPrivilegesDAO {
	dao := &RolesPrivilegesDAO{
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[model.RolesPrivilegesPrimaryKey, dbmodel.RolesPrivileges](dao)
	dao.DAOBaseQueries = &queries
	return dao
}

func (dao *RolesPrivilegesDAO) Table() context.PostgresTable {
	return table.RolesPrivileges
}

func (dao *RolesPrivilegesDAO) InsertCols() postgres.ColumnList {
	return table.RolesPrivileges.AllColumns.Except(
		table.RolesPrivileges.CreatedAt,
	)
}

func (dao *RolesPrivilegesDAO) UpdateCols() postgres.ColumnList {
	return table.RolesPrivileges.AllColumns.Except(
		table.RolesPrivileges.CreatedAt,
	)
}

func (dao *RolesPrivilegesDAO) AllCols() postgres.ColumnList {
	return table.RolesPrivileges.AllColumns
}

func (dao *RolesPrivilegesDAO) OnConflictCols() postgres.ColumnList {
	return []postgres.Column{
		table.RolesPrivileges.PrivilegeID,
		table.RolesPrivileges.RoleID,
	}
}

func (dao *RolesPrivilegesDAO) UpdateOnConflictCols() []postgres.ColumnAssigment {
	return []postgres.ColumnAssigment{
		table.RolesPrivileges.PrivilegeID.SET(table.RolesPrivileges.PrivilegeID),
		table.RolesPrivileges.RoleID.SET(table.RolesPrivileges.RoleID),
	}
}

func (dao *RolesPrivilegesDAO) PKMatch(pk model.RolesPrivilegesPrimaryKey) postgres.BoolExpression {
	return table.RolesPrivileges.
		PrivilegeID.EQ(postgres.Int(pk.PrivilegeID)).
		AND(table.RolesPrivileges.RoleID.EQ(postgres.Int16(pk.RoleID)))
}

func (dao *RolesPrivilegesDAO) GetUpdatedAt(row *dbmodel.RolesPrivileges) *time.Time {
	return nil
}
