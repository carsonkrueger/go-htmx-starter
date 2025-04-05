package database

import (
	"database/sql"

	"github.com/carsonkrueger/main/database/authDAO"
	"github.com/carsonkrueger/main/interfaces"
)

type daoManager struct {
	usersDAO                      interfaces.IUsersDAO
	privilegesDAO                 interfaces.IPrivilegeDAO
	privilegesLevelsDAO           interfaces.IPrivilegeLevelsDAO
	sessionsDAO                   interfaces.ISessionsDAO
	privilegesLevelsPrivilegesDAO interfaces.IPrivilegeLevelsPrivilegesDAO
	db                            *sql.DB
}

func NewDAOManager(db *sql.DB) interfaces.IDAOManager {
	return &daoManager{
		db: db,
	}
}

func (dm *daoManager) UsersDAO() interfaces.IUsersDAO {
	if dm.usersDAO == nil {
		dm.usersDAO = authDAO.NewUsersDAO(dm.db)
	}
	return dm.usersDAO
}

func (dm *daoManager) PrivilegeDAO() interfaces.IPrivilegeDAO {
	if dm.privilegesDAO == nil {
		dm.privilegesDAO = authDAO.NewPrivilegesDAO(dm.db)
	}
	return dm.privilegesDAO
}

func (dm *daoManager) SessionsDAO() interfaces.ISessionsDAO {
	if dm.sessionsDAO == nil {
		dm.sessionsDAO = authDAO.NewSessionsDAO(dm.db)
	}
	return dm.sessionsDAO
}

func (dm *daoManager) PrivilegeLevelsDAO() interfaces.IPrivilegeLevelsDAO {
	if dm.privilegesLevelsDAO == nil {
		dm.privilegesLevelsDAO = authDAO.NewPrivilegeLevelsDAO(dm.db)
	}
	return dm.privilegesLevelsDAO
}

func (dm *daoManager) PrivilegeLevelsPrivilegesDAO() interfaces.IPrivilegeLevelsPrivilegesDAO {
	if dm.privilegesLevelsPrivilegesDAO == nil {
		dm.privilegesLevelsPrivilegesDAO = authDAO.NewPrivilegeLevelsPrivilegesDAO(dm.db)
	}
	return dm.privilegesLevelsPrivilegesDAO
}
