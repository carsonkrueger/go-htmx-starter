package dao

import (
	"time"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/carsonkrueger/main/models/auth_models"
	"github.com/go-jet/jet/v2/postgres"
)

type privilegeLevelsPrivilegesDAO struct {
	context.DAOBaseQueries[auth_models.PrivilegeLevelsPrivilegesPrimaryKey, model.PrivilegeLevelsPrivileges]
}

func NewPrivilegeLevelsPrivilegesDAO() *privilegeLevelsPrivilegesDAO {
	dao := &privilegeLevelsPrivilegesDAO{
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[auth_models.PrivilegeLevelsPrivilegesPrimaryKey, model.PrivilegeLevelsPrivileges](dao)
	dao.DAOBaseQueries = &queries
	return dao
}

func (dao *privilegeLevelsPrivilegesDAO) Table() context.PostgresTable {
	return table.PrivilegeLevelsPrivileges
}

func (dao *privilegeLevelsPrivilegesDAO) InsertCols() postgres.ColumnList {
	return table.PrivilegeLevelsPrivileges.AllColumns.Except(
		table.PrivilegeLevelsPrivileges.CreatedAt,
		table.PrivilegeLevelsPrivileges.UpdatedAt,
	)
}

func (dao *privilegeLevelsPrivilegesDAO) UpdateCols() postgres.ColumnList {
	return table.PrivilegeLevelsPrivileges.AllColumns.Except(
		table.PrivilegeLevelsPrivileges.CreatedAt,
	)
}

func (dao *privilegeLevelsPrivilegesDAO) AllCols() postgres.ColumnList {
	return table.PrivilegeLevelsPrivileges.AllColumns
}

func (dao *privilegeLevelsPrivilegesDAO) OnConflictCols() postgres.ColumnList {
	return []postgres.Column{
		table.PrivilegeLevelsPrivileges.PrivilegeID,
		table.PrivilegeLevelsPrivileges.PrivilegeLevelID,
	}
}

func (dao *privilegeLevelsPrivilegesDAO) UpdateOnConflictCols() []postgres.ColumnAssigment {
	return []postgres.ColumnAssigment{
		table.PrivilegeLevelsPrivileges.PrivilegeID.SET(table.PrivilegeLevelsPrivileges.PrivilegeID),
		table.PrivilegeLevelsPrivileges.PrivilegeLevelID.SET(table.PrivilegeLevelsPrivileges.PrivilegeLevelID),
	}
}

func (dao *privilegeLevelsPrivilegesDAO) PKMatch(pk auth_models.PrivilegeLevelsPrivilegesPrimaryKey) postgres.BoolExpression {
	return table.PrivilegeLevelsPrivileges.
		PrivilegeID.EQ(postgres.Int(pk.PrivilegeID)).
		AND(table.PrivilegeLevelsPrivileges.PrivilegeLevelID.EQ(postgres.Int(pk.PrivilegeLevelID)))
}

func (dao *privilegeLevelsPrivilegesDAO) GetUpdatedAt(row *model.PrivilegeLevelsPrivileges) *time.Time {
	return row.UpdatedAt
}
