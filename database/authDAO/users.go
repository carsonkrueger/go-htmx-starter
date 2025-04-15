package authDAO

import (
	"database/sql"
	"time"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/carsonkrueger/main/tools"
	"github.com/go-jet/jet/v2/postgres"
)

type usersDAO struct {
	db *sql.DB
}

func NewUsersDAO(db *sql.DB) interfaces.IUsersDAO {
	return &usersDAO{
		db: db,
	}
}

func (dao *usersDAO) GetById(id int64) (*model.Users, error) {
	var user model.Users
	err := table.Users.SELECT(table.Users.AllColumns).
		FROM(table.Users).
		WHERE(table.Users.ID.EQ(postgres.Int(id))).
		LIMIT(1).
		Query(dao.db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *usersDAO) Insert(row *model.Users) error {
	return table.Users.
		INSERT(table.Users.AllColumns.Except(table.Users.ID, table.Users.CreatedAt, table.Users.UpdatedAt)).
		MODEL(row).
		RETURNING(table.Users.ID).
		Query(dao.db, row)
}

func (dao *usersDAO) InsertMany(rows []*model.Users) error {
	if len(rows) == 0 {
		return nil
	}
	return table.Users.
		INSERT(table.Users.AllColumns.Except(table.Users.ID, table.Users.CreatedAt, table.Users.UpdatedAt)).
		MODELS(rows).
		RETURNING(table.Users.ID).
		Query(dao.db, &rows)
}

func (dao *usersDAO) Upsert(row *model.Users, colsUpdate ...postgres.ColumnAssigment) error {
	if len(colsUpdate) == 0 {
		colsUpdate = []postgres.ColumnAssigment{
			table.Users.Email.SET(postgres.String(row.Email)),
			table.Users.FirstName.SET(postgres.String(row.FirstName)),
			table.Users.LastName.SET(postgres.String(row.LastName)),
		}
	}
	row.UpdatedAt = tools.Ptr(time.Now())

	return table.Users.
		INSERT(table.Users.AllColumns.Except(table.Users.ID, table.Users.CreatedAt)).
		MODEL(row).
		ON_CONFLICT(table.Users.ID, table.Users.Email).
		DO_UPDATE(postgres.SET(colsUpdate...)).
		RETURNING(table.Users.ID).
		Query(dao.db, row)
}

func (dao *usersDAO) UpsertMany(rows []*model.Users, colsUpdate ...postgres.ColumnAssigment) error {
	if len(colsUpdate) == 0 {
		colsUpdate = []postgres.ColumnAssigment{
			table.Users.Email.SET(table.Users.Email),
			table.Users.FirstName.SET(table.Users.FirstName),
			table.Users.LastName.SET(table.Users.LastName),
			table.Users.UpdatedAt.SET(table.Users.UpdatedAt),
		}
	}

	now := time.Now()
	for _, r := range rows {
		r.UpdatedAt = &now
	}

	return table.Users.
		INSERT(table.Users.AllColumns.Except(table.Users.ID, table.Users.CreatedAt)).
		MODELS(rows).
		ON_CONFLICT(table.Users.Email).
		DO_UPDATE(postgres.SET(colsUpdate...)).
		RETURNING(table.Privileges.ID).
		Query(dao.db, &rows)
}

func (dao *usersDAO) Update(row *model.Users) error {
	row.UpdatedAt = tools.Ptr(time.Now())
	_, err := table.Users.
		UPDATE(table.Users.AllColumns.Except(table.Users.ID, table.Users.CreatedAt)).
		MODEL(row).
		WHERE(table.Users.ID.EQ(postgres.Int(row.ID))).
		Exec(dao.db)
	return err
}

func (dao *usersDAO) Delete(id int64) error {
	_, err := table.Users.
		DELETE().
		WHERE(table.Users.ID.EQ(postgres.Int(id))).
		Exec(dao.db)
	return err
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

func (dao *usersDAO) GetAll() (*[]model.Users, error) {
	var rows []model.Users
	err := table.Users.
		SELECT(table.Users.AllColumns).
		Query(dao.db, &rows)
	if err != nil {
		return nil, err
	}
	return &rows, nil
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
