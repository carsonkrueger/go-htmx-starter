package dao

import (
	gctx "context"
	"time"

	"github.com/carsonkrueger/main/internal/constant"
	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/database/gen/go_starter_db/auth/table"
	dbmodel "github.com/carsonkrueger/main/pkg/db/auth/model"
	"github.com/carsonkrueger/main/pkg/model"
	"github.com/go-jet/jet/v2/postgres"
)

type privilegesDAO struct {
	context.DAOBaseQueries[int64, dbmodel.Privileges]
}

func NewPrivilegesDAO() *privilegesDAO {
	dao := &privilegesDAO{
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[int64, dbmodel.Privileges](dao)
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

func (dao *privilegesDAO) GetUpdatedAt(row *dbmodel.Privileges) *time.Time {
	return row.UpdatedAt
}

func (dao *privilegesDAO) GetAllJoined(ctx gctx.Context) ([]model.RolesPrivilegeJoin, error) {
	var res []model.RolesPrivilegeJoin

	err := table.Roles.
		SELECT(
			table.Roles.AllColumns,
			table.Privileges.AllColumns,
		).
		FROM(
			table.RolesPrivileges.
				LEFT_JOIN(table.Roles, table.Roles.ID.EQ(table.RolesPrivileges.RoleID)).
				LEFT_JOIN(table.Privileges, table.Privileges.ID.EQ(table.RolesPrivileges.PrivilegeID)),
		).
		ORDER_BY(table.Roles.Name.ASC(), table.Privileges.Name.ASC()).
		Query(context.GetDB(ctx), &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (dao *privilegesDAO) GetPrivilegesByRoleID(ctx gctx.Context, roleID int64) ([]dbmodel.Roles, error) {
	var privileges []dbmodel.Roles
	err := table.RolesPrivileges.
		SELECT(
			table.RolesPrivileges.RoleID,
			table.RolesPrivileges.PrivilegeID,
			table.Privileges.AllColumns,
		).
		FROM(
			table.RolesPrivileges.
				LEFT_JOIN(table.Privileges, table.Privileges.ID.EQ(table.RolesPrivileges.RoleID)),
		).
		Query(context.GetDB(ctx), &privileges)
	if err != nil {
		return privileges, err
	}
	return privileges, nil
}

func (dao *privilegesDAO) GetManyByName(ctx gctx.Context, names []constant.PrivilegeName) ([]dbmodel.Privileges, error) {
	var privileges []dbmodel.Privileges

	// Handle empty slice case
	if len(names) == 0 {
		return privileges, nil
	}

	exprs := make([]postgres.Expression, len(names))
	for i, name := range names {
		exprs[i] = postgres.String(string(name))
	}

	err := table.Privileges.
		SELECT(table.Privileges.AllColumns).
		WHERE(table.Privileges.Name.IN(exprs...)).
		QueryContext(ctx, context.GetDB(ctx), &privileges)
	if err != nil {
		return nil, err
	}

	return privileges, nil
}
