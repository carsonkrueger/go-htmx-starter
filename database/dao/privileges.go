package dao

import (
	gctx "context"
	"time"

	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/carsonkrueger/main/models/auth_models"
	"github.com/go-jet/jet/v2/postgres"
)

type privilegesDAO struct {
	context.DAOBaseQueries[int64, model.Privileges]
}

func NewPrivilegesDAO() *privilegesDAO {
	dao := &privilegesDAO{
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[int64, model.Privileges](dao)
	dao.DAOBaseQueries = &queries
	return dao
}

func (dao *privilegesDAO) Table() context.PostgresTable {
	return table.Privileges
}

func (dao *privilegesDAO) InsertCols() postgres.ColumnList {
	return table.Privileges.AllColumns.Except(
		table.Privileges.ID,
		table.Privileges.CreatedAt,
		table.Privileges.UpdatedAt,
	)
}

func (dao *privilegesDAO) UpdateCols() postgres.ColumnList {
	return table.Privileges.AllColumns.Except(
		table.Privileges.ID,
		table.Privileges.CreatedAt,
	)
}

func (dao *privilegesDAO) AllCols() postgres.ColumnList {
	return table.Privileges.AllColumns
}

func (dao *privilegesDAO) OnConflictCols() postgres.ColumnList {
	return []postgres.Column{table.Privileges.Name}
}

func (dao *privilegesDAO) UpdateOnConflictCols() []postgres.ColumnAssigment {
	return []postgres.ColumnAssigment{
		table.Privileges.Name.SET(table.Privileges.EXCLUDED.Name),
	}
}

func (dao *privilegesDAO) PKMatch(pk int64) postgres.BoolExpression {
	return table.Privileges.ID.EQ(postgres.Int(pk))
}

func (dao *privilegesDAO) GetUpdatedAt(row *model.Privileges) *time.Time {
	return row.UpdatedAt
}

func (dao *privilegesDAO) GetAllJoined(ctx gctx.Context) ([]auth_models.JoinedPrivilegesRaw, error) {
	var res []auth_models.JoinedPrivilegesRaw

	err := table.PrivilegeLevels.
		SELECT(
			table.PrivilegeLevelsPrivileges.PrivilegeLevelID.AS("JoinedPrivilegesRaw.LevelID"),
			table.PrivilegeLevels.Name.AS("JoinedPrivilegesRaw.LevelName"),
			table.PrivilegeLevelsPrivileges.PrivilegeID.AS("JoinedPrivilegesRaw.PrivilegeID"),
			table.Privileges.Name.AS("JoinedPrivilegesRaw.PrivilegeName"),
			table.Privileges.CreatedAt.AS("JoinedPrivilegesRaw.PrivilegeCreatedAt"),
		).
		FROM(
			table.PrivilegeLevelsPrivileges.
				LEFT_JOIN(table.PrivilegeLevels, table.PrivilegeLevels.ID.EQ(table.PrivilegeLevelsPrivileges.PrivilegeLevelID)).
				LEFT_JOIN(table.Privileges, table.Privileges.ID.EQ(table.PrivilegeLevelsPrivileges.PrivilegeID)),
		).
		ORDER_BY(table.PrivilegeLevels.Name.ASC(), table.Privileges.Name.ASC()).
		Query(context.GetDB(ctx), &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (dao *privilegesDAO) GetPrivilegesByLevelID(ctx gctx.Context, levelID int64) ([]model.PrivilegeLevels, error) {
	var privileges []model.PrivilegeLevels
	err := table.PrivilegeLevelsPrivileges.
		SELECT(
			table.PrivilegeLevelsPrivileges.PrivilegeLevelID,
			table.PrivilegeLevelsPrivileges.PrivilegeID,
			table.Privileges.AllColumns,
		).
		FROM(
			table.PrivilegeLevelsPrivileges.
				LEFT_JOIN(table.Privileges, table.Privileges.ID.EQ(table.PrivilegeLevelsPrivileges.PrivilegeLevelID)),
		).
		Query(context.GetDB(ctx), &privileges)
	if err != nil {
		return privileges, err
	}
	return privileges, nil
}
