package services

import (
	gctx "context"
	"fmt"

	"github.com/carsonkrueger/main/internal/context"
	dbmodel "github.com/carsonkrueger/main/pkg/db/auth/model"
	"github.com/carsonkrueger/main/pkg/model"
)

type privilegesService struct {
	*context.AppContext
}

func NewPrivilegesService(ctx *context.AppContext) *privilegesService {
	return &privilegesService{ctx}
}

func (ps *privilegesService) CreatePrivilegeAssociation(ctx gctx.Context, roleID int16, privID int64) error {
	joinRow := dbmodel.RolesPrivileges{
		RoleID:      roleID,
		PrivilegeID: privID,
	}
	rolesPrivsDAO := ps.DM().RolesPrivilegesDAO()

	if err := rolesPrivsDAO.Upsert(ctx, &joinRow); err != nil {
		return fmt.Errorf("failed to associate role to privilege %d-%d: %w", roleID, privID, err)
	}

	return nil
}

func (ps *privilegesService) CreateRole(ctx gctx.Context, name string) error {
	row := dbmodel.Roles{
		Name: name,
	}
	rolesDAO := ps.DM().RolesDAO()
	if err := rolesDAO.Insert(ctx, &row); err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	return nil
}

func (ps *privilegesService) HasPermissionsByIDS(ctx gctx.Context, roleID int16, privIDs []int64) bool {
	pks := make([]model.RolesPrivilegesPrimaryKey, len(privIDs))
	for i, privID := range privIDs {
		pks[i] = model.RolesPrivilegesPrimaryKey{
			PrivilegeID: privID,
			RoleID:      roleID,
		}
	}
	rows, err := ps.DM().RolesPrivilegesDAO().GetMany(ctx, pks)
	return len(rows) > 0 && err == nil
}

func (ps *privilegesService) DeletePrivilegeAssociation(ctx gctx.Context, roleID int16, privID int64) error {
	pk := model.RolesPrivilegesPrimaryKey{
		PrivilegeID: privID,
		RoleID:      roleID,
	}
	if err := ps.DM().RolesPrivilegesDAO().Delete(ctx, pk); err != nil {
		return fmt.Errorf("failed to delete privilege association: %w", err)
	}

	return nil
}

func (us *privilegesService) SetUserRole(ctx gctx.Context, roleID int16, userID int64) error {
	roleDAO := us.DM().RolesDAO()
	role, err := roleDAO.GetOne(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to fetch role %d: %w", roleID, err)
	}

	userDAO := us.DM().UsersDAO()
	user, err := userDAO.GetOne(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user %d: %w", userID, err)
	}

	user.RoleID = role.ID
	if err := userDAO.Update(ctx, &user, user.ID); err != nil {
		return fmt.Errorf("failed to update user %d: %w", userID, err)
	}

	return nil
}
