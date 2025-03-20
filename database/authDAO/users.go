package authDAO

import (
	"database/sql"
	"time"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/carsonkrueger/main/interfaces"
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

func (dao *usersDAO) Insert(row *model.Users) (int64, error) {
	err := table.Users.
		INSERT(table.Users.Email, table.Users.FirstName, table.Users.LastName, table.Users.Password, table.Users.AuthToken, table.Users.AuthTokenCreatedAt, table.Users.PrivilegeLevelID).
		VALUES(row.Email, row.FirstName, row.LastName, row.Password, row.AuthToken, postgres.TimestampT(time.Now()), row.PrivilegeLevelID).
		RETURNING(table.Users.ID).
		Query(dao.db, row)

	if err != nil {
		return -1, err
	}
	return row.ID, nil
}

// Returns ID int64 if inserted.
// Parameter cols_update are the columns to be updated on conflict - if not provided, a few columns are updated
func (dao *usersDAO) Upsert(row *model.Users, colsUpdate ...postgres.ColumnAssigment) (int64, error) {
	if len(colsUpdate) == 0 {
		colsUpdate = []postgres.ColumnAssigment{
			table.Users.Email.SET(postgres.String(row.Email)),
			table.Users.FirstName.SET(postgres.String(row.FirstName)),
			table.Users.LastName.SET(postgres.String(row.LastName)),
			table.Users.UpdatedAt.SET(postgres.TimestampT(time.Now())),
		}
	}

	res, err := table.Users.
		INSERT(table.Users.AllColumns.Except(table.Users.ID, table.Users.CreatedAt)).
		MODEL(row).
		ON_CONFLICT(table.Users.ID, table.Users.Email).
		DO_UPDATE(postgres.SET(colsUpdate...)).
		RETURNING(table.Users.ID).
		Exec(dao.db)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	// if err then row was not inserted
	if err != nil {
		id = -1
	}

	return id, nil
}

func (dao *usersDAO) Update(row *model.Users) error {
	_, err := table.Users.
		UPDATE(table.Users.AllColumns).
		MODEL(row).
		WHERE(table.Users.ID.EQ(postgres.Int(row.ID))).
		Exec(dao.db)
	if err != nil {
		return err
	}
	return nil
}

func (dao *usersDAO) Delete(id int64) error {
	_, err := table.Users.
		DELETE().
		WHERE(table.Users.ID.EQ(postgres.Int(id))).
		Exec(dao.db)
	if err != nil {
		return err
	}
	return nil
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

func (dao *usersDAO) GetPrivilegeLevelID(userID int64, token string) (int64, error) {
	var res PrivilegeLevelIDResponse

	err := table.Users.
		SELECT(table.Users.PrivilegeLevelID.AS("PrivilegeLevelIDResponse.PrivilegeID")).
		FROM(table.Users).
		WHERE(table.Users.ID.EQ(postgres.Int(userID)).
			AND(table.Users.AuthToken.EQ(postgres.String(token)))).
		LIMIT(1).
		Query(dao.db, &res)

	if err != nil {
		return -1, err
	}
	return res.PrivilegeID, nil
}

func (dao *usersDAO) UpdateAuthToken(id int64, authToken string) error {
	_, err := table.Users.
		UPDATE(table.Users.AuthToken).
		SET(authToken).
		WHERE(table.Users.ID.EQ(postgres.Int(id))).
		Exec(dao.db)
	if err != nil {
		return err
	}
	return nil
}

func (dao *usersDAO) GetAll() ([]*model.Users, error) {
	// not implemented
	return []*model.Users{}, nil
}
