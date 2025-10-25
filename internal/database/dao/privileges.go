package dao

import (
	gctx "context"
	"time"

	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/gen/go_starter_db/auth/table"
	"github.com/carsonkrueger/main/pkg/model"
	"github.com/carsonkrueger/main/pkg/model/db/auth"
	"github.com/go-jet/jet/v2/postgres"
)

type privilegesDAO struct {
	context.DAOBaseQueries[int64, auth.Privileges]
}

func NewPrivilegesDAO() *privilegesDAO {
	dao := &privilegesDAO{
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[int64, auth.Privileges](dao)
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

func (dao *privilegesDAO) GetUpdatedAt(row *auth.Privileges) *time.Time {
	return row.UpdatedAt
}

func (dao *privilegesDAO) GetAllJoined(ctx gctx.Context) ([]model.JoinedPrivilegesRaw, error) {
	var res []model.JoinedPrivilegesRaw

	err := table.Roles.
		SELECT(
			table.RolesPrivileges.RoleID.AS("JoinedPrivilegesRaw.RoleID"),
			table.Roles.Name.AS("JoinedPrivilegesRaw.RoleName"),
			table.RolesPrivileges.PrivilegeID.AS("JoinedPrivilegesRaw.PrivilegeID"),
			table.Privileges.Name.AS("JoinedPrivilegesRaw.PrivilegeName"),
			table.Privileges.CreatedAt.AS("JoinedPrivilegesRaw.PrivilegeCreatedAt"),
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

func (dao *privilegesDAO) GetPrivilegesByRoleID(ctx gctx.Context, roleID int64) ([]auth.Roles, error) {
	var privileges []auth.Roles
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

func (dao *privilegesDAO) GetManyByName(ctx gctx.Context, names []string) ([]auth.Privileges, error) {
	var privileges []auth.Privileges

	// Handle empty slice case
	if len(names) == 0 {
		return privileges, nil
	}

	exprs := make([]postgres.Expression, len(names))
	for i, name := range names {
		exprs[i] = postgres.String(name)
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
