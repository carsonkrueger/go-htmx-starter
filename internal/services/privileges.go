package services

import (
	gctx "context"
	"errors"
	"strconv"

	"github.com/carsonkrueger/main/internal/context"
	dbmodel "github.com/carsonkrueger/main/pkg/db/auth/model"
	"github.com/carsonkrueger/main/pkg/model"
	"go.uber.org/zap"
)

type privilegesService struct {
	*context.AppContext
}

func NewPrivilegesService(ctx *context.AppContext) *privilegesService {
	return &privilegesService{ctx}
}

func (ps *privilegesService) CreatePrivilegeAssociation(ctx gctx.Context, roleID int16, privID int64) error {
	lgr := ps.Lgr("AddPermission")
	lgr.Info("Role:Privilege", zap.String(strconv.FormatInt(int64(roleID), 10), strconv.FormatInt(int64(privID), 10)))

	joinRow := dbmodel.RolesPrivileges{
		RoleID:      roleID,
		PrivilegeID: privID,
	}
	rolesPrivsDAO := ps.DM().RolesPrivilegesDAO()

	if err := rolesPrivsDAO.Upsert(ctx, &joinRow); err != nil {
		lgr.Error("Failed to insert roles privileges", zap.Error(err), zap.Any("rolesPrivileges", joinRow))
		return err
	}

	return nil
}

func (ps *privilegesService) CreateRole(ctx gctx.Context, name string) error {
	lgr := ps.Lgr("CreateRole")
	lgr.Info("Called")

	row := dbmodel.Roles{
		Name: name,
	}
	rolesDAO := ps.DM().RolesDAO()
	if err := rolesDAO.Insert(ctx, &row); err != nil {
		lgr.Error("Failed to create role", zap.Error(err))
		return errors.New("Failed to create role")
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
	lgr := ps.Lgr("DeletePrivilegeAssociation")
	lgr.Info("Called")

	pk := model.RolesPrivilegesPrimaryKey{
		PrivilegeID: privID,
		RoleID:      roleID,
	}
	if err := ps.DM().RolesPrivilegesDAO().Delete(ctx, pk); err != nil {
		return nil
	}

	return nil
}

func (us *privilegesService) SetUserRole(ctx gctx.Context, roleID int16, userID int64) error {
	lgr := us.Lgr("SetUserRole")
	lgr.Info("Called")

	roleDAO := us.DM().RolesDAO()
	role, err := roleDAO.GetOne(ctx, roleID)
	if err != nil {
		lgr.Error("Error fetching role", zap.Error(err))
		return err
	}

	userDAO := us.DM().UsersDAO()
	user, err := userDAO.GetOne(ctx, userID)
	if err != nil {
		lgr.Error("Error fetching user", zap.Error(err))
		return err
	}

	user.RoleID = role.ID
	if err := userDAO.Update(ctx, &user, user.ID); err != nil {
		lgr.Error("Error updating user", zap.Error(err))
		return err
	}

	return nil
}
