package DAO

import (
	"database/sql"
	"time"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/go-jet/jet/v2/postgres"
)

type usersDAO struct {
	db *sql.DB
	DAOBaseQueries[int64, model.Users]
}

func newUsersDAO(db *sql.DB) UsersDAO {
	dao := &usersDAO{
		db:             db,
		DAOBaseQueries: nil,
	}
	queries := newDAOQueryable[int64, model.Users](dao)
	dao.DAOBaseQueries = &queries
	return dao
}

func (dao *usersDAO) Table() PostgresTable {
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

func (dao *usersDAO) GetUpdatedAt(row *model.Users) *time.Time {
	return row.UpdatedAt
}

func (dao *usersDAO) GetByEmail(email string) (*model.Users, error) {
	var user model.Users
	err := table.Users.
		SELECT(table.Users.AllColumns).
		FROM(table.Users).
		WHERE(table.Users.Email.EQ(postgres.String(email))).
		LIMIT(1).
		Query(dao.db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type PrivilegeLevelIDResponse struct {
	PrivilegeID int64
}

func (dao *usersDAO) GetPrivilegeLevelID(userID int64) (*int64, error) {
	var res PrivilegeLevelIDResponse

	err := table.Users.
		SELECT(table.Users.PrivilegeLevelID.AS("PrivilegeLevelIDResponse.PrivilegeID")).
		FROM(table.Users).
		WHERE(table.Users.ID.EQ(postgres.Int(userID))).
		LIMIT(1).
		Query(dao.db, &res)

	if err != nil {
		return nil, err
	}
	return &res.PrivilegeID, nil
}

func (dao *usersDAO) GetUserPrivilegeJoinAll() (*[]authModels.UserPrivilegeLevelJoin, error) {
	var rows []authModels.UserPrivilegeLevelJoin
	err := table.Users.
		LEFT_JOIN(table.PrivilegeLevels, table.Users.PrivilegeLevelID.EQ(table.PrivilegeLevels.ID)).
		SELECT(table.Users.AllColumns, table.PrivilegeLevels.Name).
		Query(dao.db, &rows)
	if err != nil {
		return nil, err
	}
	return &rows, nil
}
