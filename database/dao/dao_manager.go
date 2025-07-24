package dao

import (
	"github.com/carsonkrueger/main/context"
)

type daoManager struct {
	usersDAO                      context.UsersDAO
	privilegesDAO                 context.PrivilegeDAO
	privilegesLevelsDAO           context.PrivilegeLevelsDAO
	sessionsDAO                   context.SessionsDAO
	privilegesLevelsPrivilegesDAO context.PrivilegeLevelsPrivilegesDAO
	// INSERT DAO
}

func NewDAOManager() context.DAOManager {
	return &daoManager{}
}

func (dm *daoManager) UsersDAO() context.UsersDAO {
	if dm.usersDAO == nil {
		dm.usersDAO = NewUsersDAO()
	}
	return dm.usersDAO
}

func (dm *daoManager) PrivilegeDAO() context.PrivilegeDAO {
	if dm.privilegesDAO == nil {
		dm.privilegesDAO = NewPrivilegesDAO()
	}
	return dm.privilegesDAO
}

func (dm *daoManager) SessionsDAO() context.SessionsDAO {
	if dm.sessionsDAO == nil {
		dm.sessionsDAO = NewSessionsDAO()
	}
	return dm.sessionsDAO
}

func (dm *daoManager) PrivilegeLevelsDAO() context.PrivilegeLevelsDAO {
	if dm.privilegesLevelsDAO == nil {
		dm.privilegesLevelsDAO = NewPrivilegeLevelsDAO()
	}
	return dm.privilegesLevelsDAO
}

func (dm *daoManager) PrivilegeLevelsPrivilegesDAO() context.PrivilegeLevelsPrivilegesDAO {
	if dm.privilegesLevelsPrivilegesDAO == nil {
		dm.privilegesLevelsPrivilegesDAO = NewPrivilegeLevelsPrivilegesDAO()
	}
	return dm.privilegesLevelsPrivilegesDAO
}

// INSERT INIT DAO
