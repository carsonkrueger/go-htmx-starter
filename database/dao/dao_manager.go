package dao

import (
	"github.com/carsonkrueger/main/context"

	"database/sql"
)

type daoManager struct {
	db                            *sql.DB
	usersDAO                      context.UsersDAO
	privilegesDAO                 context.PrivilegeDAO
	privilegesLevelsDAO           context.PrivilegeLevelsDAO
	sessionsDAO                   context.SessionsDAO
	privilegesLevelsPrivilegesDAO context.PrivilegeLevelsPrivilegesDAO
	// INSERT DAO
}

func NewDAOManager(db *sql.DB) context.DAOManager {
	return &daoManager{
		db: db,
	}
}

func (dm *daoManager) UsersDAO() context.UsersDAO {
	if dm.usersDAO == nil {
		dm.usersDAO = NewUsersDAO(dm.db)
	}
	return dm.usersDAO
}

func (dm *daoManager) PrivilegeDAO() context.PrivilegeDAO {
	if dm.privilegesDAO == nil {
		dm.privilegesDAO = NewPrivilegesDAO(dm.db)
	}
	return dm.privilegesDAO
}

func (dm *daoManager) SessionsDAO() context.SessionsDAO {
	if dm.sessionsDAO == nil {
		dm.sessionsDAO = NewSessionsDAO(dm.db)
	}
	return dm.sessionsDAO
}

func (dm *daoManager) PrivilegeLevelsDAO() context.PrivilegeLevelsDAO {
	if dm.privilegesLevelsDAO == nil {
		dm.privilegesLevelsDAO = NewPrivilegeLevelsDAO(dm.db)
	}
	return dm.privilegesLevelsDAO
}

func (dm *daoManager) PrivilegeLevelsPrivilegesDAO() context.PrivilegeLevelsPrivilegesDAO {
	if dm.privilegesLevelsPrivilegesDAO == nil {
		dm.privilegesLevelsPrivilegesDAO = NewPrivilegeLevelsPrivilegesDAO(dm.db)
	}
	return dm.privilegesLevelsPrivilegesDAO
}

// INSERT INIT DAO
