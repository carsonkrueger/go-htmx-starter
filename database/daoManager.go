package database

import (
	"database/sql"

	"github.com/carsonkrueger/main/database/auth"
	"github.com/carsonkrueger/main/interfaces"
)

type daoManager struct {
	usersDAO      interfaces.IUsersDAO
	privilegesDAO interfaces.IPrivilegeDAO
	db            *sql.DB
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

func (dm *daoManager) PrivilegeDAO() interfaces.IPrivilegeDAO {
	if dm.privilegesDAO == nil {
		dm.privilegesDAO = auth.NewPrivilegesDAO(dm.db)
	}
	return dm.privilegesDAO
}
