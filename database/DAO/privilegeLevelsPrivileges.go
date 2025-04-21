package DAO

import (
	"database/sql"
	"time"

	"github.com/carsonkrueger/main/database"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/go-jet/jet/v2/postgres"
)

type privilegeLevelsPrivilegesDAO struct {
	db *sql.DB
	interfaces.IDAOBaseQueries[authModels.PrivilegeLevelsPrivilegesPrimaryKey, model.PrivilegeLevelsPrivileges]
}

func newPrivilegeLevelsPrivilegesDAO(db *sql.DB) *privilegeLevelsPrivilegesDAO {
	dao := &privilegeLevelsPrivilegesDAO{
		db:              db,
		IDAOBaseQueries: nil,
	}
	queries := database.NewDAOQueryable(dao)
	dao.IDAOBaseQueries = &queries
	return dao
}

func (dao *privilegeLevelsPrivilegesDAO) Table() interfaces.IPostgresTable {
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
		table.PrivilegeLevelsPrivileges.PrivilegeID,
		table.PrivilegeLevelsPrivileges.PrivilegeLevelID,
	)
}

func (dao *privilegeLevelsPrivilegesDAO) AllCols() postgres.ColumnList {
	return table.PrivilegeLevelsPrivileges.AllColumns
}

func (dao *privilegeLevelsPrivilegesDAO) OnConflictCols() postgres.ColumnList {
	return []postgres.Column{}
}

func (dao *privilegeLevelsPrivilegesDAO) UpdateOnConflictCols() []postgres.ColumnAssigment {
	return []postgres.ColumnAssigment{}
}

func (dao *privilegeLevelsPrivilegesDAO) PKMatch(pk authModels.PrivilegeLevelsPrivilegesPrimaryKey) postgres.BoolExpression {
	return table.PrivilegeLevelsPrivileges.
		PrivilegeID.EQ(postgres.Int(pk.PrivilegeID)).
		AND(table.PrivilegeLevelsPrivileges.PrivilegeLevelID.EQ(postgres.Int(pk.PrivilegeLevelID)))
}

func (dao *privilegeLevelsPrivilegesDAO) GetUpdatedAt(row *model.PrivilegeLevelsPrivileges) *time.Time {
	return row.UpdatedAt
}
