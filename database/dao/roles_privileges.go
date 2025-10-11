package dao

import (
	"time"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/gen/go_starter_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_starter_db/auth/table"
	"github.com/carsonkrueger/main/models/auth_models"
	"github.com/go-jet/jet/v2/postgres"
)

type RolesPrivilegesDAO struct {
	context.DAOBaseQueries[auth_models.RolesPrivilegesPrimaryKey, model.RolesPrivileges]
}

func NewRolesPrivilegesDAO() *RolesPrivilegesDAO {
	dao := &RolesPrivilegesDAO{
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[auth_models.RolesPrivilegesPrimaryKey, model.RolesPrivileges](dao)
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

func (dao *RolesPrivilegesDAO) PKMatch(pk auth_models.RolesPrivilegesPrimaryKey) postgres.BoolExpression {
	return table.RolesPrivileges.
		PrivilegeID.EQ(postgres.Int(pk.PrivilegeID)).
		AND(table.RolesPrivileges.RoleID.EQ(postgres.Int16(pk.RoleID)))
}

func (dao *RolesPrivilegesDAO) GetUpdatedAt(row *model.RolesPrivileges) *time.Time {
	return nil
}
