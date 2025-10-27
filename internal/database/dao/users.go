package dao

import (
	gctx "context"
	"time"

	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/gen/go_starter_db/auth/table"
	dbmodel "github.com/carsonkrueger/main/pkg/db/auth/model"
	"github.com/carsonkrueger/main/pkg/model"
	"github.com/go-jet/jet/v2/postgres"
)

type usersDAO struct {
	context.DAOBaseQueries[int64, dbmodel.Users]
}

func NewUsersDAO() context.UsersDAO {
	dao := &usersDAO{
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[int64, dbmodel.Users](dao)
	dao.DAOBaseQueries = &queries
	return dao
}

func (dao *usersDAO) Table() context.PostgresTable {
	return table.Users
}

func (dao *usersDAO) InsertCols() postgres.ColumnList {
	return table.Users.AllColumns.Except(
		table.Users.ID,
		table.Users.CreatedAt,
		table.Users.UpdatedAt,
	)
}

func (dao *usersDAO) UpdateCols() postgres.ColumnList {
	return table.Users.AllColumns.Except(
		table.Users.ID,
		table.Users.CreatedAt,
	)
}

func (dao *usersDAO) AllCols() postgres.ColumnList {
	return table.Users.AllColumns
}

func (dao *usersDAO) OnConflictCols() postgres.ColumnList {
	return []postgres.Column{}
}

func (dao *usersDAO) UpdateOnConflictCols() []postgres.ColumnAssigment {
	return []postgres.ColumnAssigment{}
}

func (dao *usersDAO) PKMatch(pk int64) postgres.BoolExpression {
	return table.Users.ID.EQ(postgres.Int(pk))
}

func (dao *usersDAO) GetUpdatedAt(row *dbmodel.Users) *time.Time {
	return row.UpdatedAt
}

func (dao *usersDAO) GetByEmail(ctx gctx.Context, email string) (*dbmodel.Users, error) {
	var user dbmodel.Users
	err := table.Users.
		SELECT(table.Users.AllColumns).
		FROM(table.Users).
		WHERE(table.Users.Email.EQ(postgres.String(email))).
		LIMIT(1).
		Query(context.GetDB(ctx), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type RoleIDResponse struct {
	PrivilegeID int64
}

func (dao *usersDAO) GetRoleID(ctx gctx.Context, userID int64) (*int64, error) {
	var res RoleIDResponse

	err := table.Users.
		SELECT(table.Users.RoleID.AS("RoleIDResponse.PrivilegeID")).
		FROM(table.Users).
		WHERE(table.Users.ID.EQ(postgres.Int(userID))).
		LIMIT(1).
		Query(context.GetDB(ctx), &res)

	if err != nil {
		return nil, err
	}
	return &res.PrivilegeID, nil
}

func (dao *usersDAO) GetUserPrivilegeJoinAll(ctx gctx.Context) (*[]model.UserRoleJoin, error) {
	var rows []model.UserRoleJoin
	err := table.Users.
		LEFT_JOIN(table.Roles, table.Users.RoleID.EQ(table.Roles.ID)).
		SELECT(table.Users.AllColumns, table.Roles.Name).
		Query(context.GetDB(ctx), &rows)
	if err != nil {
		return nil, err
	}
	return &rows, nil
}
