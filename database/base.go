package database

import (
	"database/sql"

	"github.com/carsonkrueger/main/database/auth"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/go-jet/jet/v2/postgres"
)

type IDAO interface {
	GetById(id int64) (*model.Users, error)
	Insert(row *model.Users) (int64, error)
	Upsert(row *model.Users, cols_update ...postgres.ColumnAssigment) (int64, error)
	Update(row *model.Users) error
	Delete(id int64) error
}

type IDAOManager interface {
	UsersDAO() auth.IUsersDAO
}

type daoManager struct {
	usersDAO auth.IUsersDAO
	db       *sql.DB
}

func NewDAOManager(db *sql.DB) IDAOManager {
	return &daoManager{
		db: db,
	}
}

func (dm *daoManager) UsersDAO() auth.IUsersDAO {
	if dm.usersDAO == nil {
		dm.usersDAO = auth.NewUsersDAO(dm.db)
	}
	return dm.usersDAO
}
