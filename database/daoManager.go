package database

import (
	"database/sql"

	"github.com/carsonkrueger/main/database/auth"
	"github.com/carsonkrueger/main/interfaces"
)

type daoManager struct {
	usersDAO interfaces.IUsersDAO
	db       *sql.DB
}

func NewDAOManager(db *sql.DB) interfaces.IDAOManager {
	return &daoManager{
		db: db,
	}
}

func (dm *daoManager) UsersDAO() interfaces.IUsersDAO {
	if dm.usersDAO == nil {
		dm.usersDAO = auth.NewUsersDAO(dm.db)
	}
	return dm.usersDAO
}
