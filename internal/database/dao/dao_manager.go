package dao

import (
	"github.com/carsonkrueger/main/internal/context"
)

type daoManager struct {
	usersDAO           context.UsersDAO
	privilegesDAO      context.PrivilegeDAO
	rolesDAO           context.RolesDAO
	sessionsDAO        context.SessionsDAO
	rolesPrivilegesDAO context.RolesPrivilegesDAO
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

func (dm *daoManager) RolesDAO() context.RolesDAO {
	if dm.rolesDAO == nil {
		dm.rolesDAO = NewRolesDAO()
	}
	return dm.rolesDAO
}

func (dm *daoManager) RolesPrivilegesDAO() context.RolesPrivilegesDAO {
	if dm.rolesPrivilegesDAO == nil {
		dm.rolesPrivilegesDAO = NewRolesPrivilegesDAO()
	}
	return dm.rolesPrivilegesDAO
}

// INSERT INIT DAO
