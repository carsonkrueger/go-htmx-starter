package services

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/carsonkrueger/main/tools"
	"github.com/go-jet/jet/v2/postgres"
)

type IUsersService interface {
	Index(id int64) (*model.Users, error)
	GetByEmail(email string) (*model.Users, error)
	Insert(row *model.Users) (int64, error)
	Upsert(row *model.Users, cols_update ...postgres.ColumnAssigment) (int64, error)
	Update(row *model.Users) error
	Delete(id int64) error
	UpdateAuthToken(id int64, authToken string) error
	Login(email string, password string) (*string, error)
	IsPermitted(userId int64, privilegeId int64) bool
}

type usersService struct {
	db *sql.DB
}

func NewUsersService(db *sql.DB) *usersService {
	return &usersService{
		db: db,
	}
}

func (us *usersService) Index(id int64) (*model.Users, error) {
	var user model.Users
	err := table.Users.SELECT(table.Users.AllColumns).FROM(table.Users).WHERE(table.Users.ID.EQ(postgres.Int(id))).FETCH_FIRST(postgres.Int(1)).ROWS_ONLY().Query(us.db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *usersService) GetByEmail(email string) (*model.Users, error) {
	var user model.Users
	err := table.Users.SELECT(table.Users.AllColumns).FROM(table.Users).WHERE(table.Users.Email.EQ(postgres.String(email))).FETCH_FIRST(postgres.Int(1)).ROWS_ONLY().Query(us.db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *usersService) Insert(row *model.Users) (int64, error) {
	err := table.Users.
		INSERT(table.Users.Email, table.Users.FirstName, table.Users.LastName, table.Users.Password, table.Users.AuthToken, table.Users.AuthTokenCreatedAt, table.Users.PrivilegeLevelID).
		VALUES(row.Email, row.FirstName, row.LastName, row.Password, row.AuthToken, postgres.TimestampT(time.Now()), row.PrivilegeLevelID).
		RETURNING(table.Users.ID).
		Query(us.db, row)

	if err != nil {
		return -1, err
	}
	return row.ID, nil
}

// Returns ID int64 if inserted.
// Parameter cols_update are the columns to be updated on conflict - if not provided, a few columns are updated
func (us *usersService) Upsert(row *model.Users, cols_update ...postgres.ColumnAssigment) (int64, error) {
	if len(cols_update) == 0 {
		cols_update = []postgres.ColumnAssigment{table.Users.Email.SET(postgres.String(row.Email)),
			table.Users.FirstName.SET(postgres.String(row.FirstName)),
			table.Users.LastName.SET(postgres.String(row.LastName)),
			table.Users.UpdatedAt.SET(postgres.TimestampT(time.Now())),
		}
	}

	res, err := table.Users.
		INSERT(table.Users.EXCLUDED.ID).
		VALUES(row).
		ON_CONFLICT(table.Users.ID, table.Users.Email).
		DO_UPDATE(postgres.SET(cols_update...)).
		RETURNING(table.Users.ID).
		Exec(us.db)
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

func (us *usersService) Update(row *model.Users) error {
	_, err := table.Users.UPDATE(table.Users.AllColumns).MODEL(row).WHERE(table.Users.ID.EQ(postgres.Int(row.ID))).Exec(us.db)
	if err != nil {
		return err
	}
	return nil
}

func (us *usersService) Delete(id int64) error {
	_, err := table.Users.DELETE().WHERE(table.Users.ID.EQ(postgres.Int(id))).Exec(us.db)
	if err != nil {
		return err
	}
	return nil
}

func (us *usersService) UpdateAuthToken(id int64, authToken string) error {
	_, err := table.Users.UPDATE(table.Users.AuthToken).SET(authToken).WHERE(table.Users.ID.EQ(postgres.Int(id))).Exec(us.db)
	if err != nil {
		return err
	}
	return nil
}

func (us *usersService) Login(email string, password string) (*string, error) {
	user, err := us.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(user.Password, "$")
	hash := tools.HashPassword(password, parts[0])

	if user.Password != hash {
		return nil, errors.New("Invalid password")
	}

	token, _ := tools.GenerateSalt()
	fullToken := fmt.Sprintf("%s$%d", token, user.ID)
	err = us.UpdateAuthToken(user.ID, fullToken)
	if err != nil {
		return nil, err
	}

	return &fullToken, nil
}

type IsPermittedResponse struct {
	PrivilegeID int64
}

func (us *usersService) IsPermitted(userId int64, privilegeId int64) bool {
	var res IsPermittedResponse

	err := table.Users.
		LEFT_JOIN(table.PrivilegeLevels, table.PrivilegeLevels.ID.EQ(table.Users.PrivilegeLevelID)).
		LEFT_JOIN(table.PrivilegeLevelsPrivileges, table.PrivilegeLevelsPrivileges.PrivilegeLevelID.EQ(table.PrivilegeLevels.ID)).
		SELECT(table.PrivilegeLevelsPrivileges.PrivilegeID.AS("IsPermittedResponse.PrivilegeID")).
		WHERE(table.PrivilegeLevelsPrivileges.PrivilegeID.EQ(postgres.Int(privilegeId)).
			AND(table.Users.ID.EQ(postgres.Int(userId)))).
		LIMIT(1).
		Query(us.db, &res)

	if err != nil {
		return false
	}

	return res.PrivilegeID == privilegeId
}
