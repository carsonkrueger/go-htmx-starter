package services

import (
	gctx "context"
	"errors"
	"fmt"
	"strconv"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/templates/ui/partials/basiclabel"
	"github.com/carsonkrueger/main/internal/templates/ui/partials/basictable"
	"github.com/carsonkrueger/main/internal/templates/ui/partials/form"
	"github.com/carsonkrueger/main/internal/templates/ui/partials/selectbox"
	dbmodel "github.com/carsonkrueger/main/pkg/db/auth/model"
	"github.com/carsonkrueger/main/pkg/model"
	"github.com/carsonkrueger/main/pkg/templui/icon"
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

func (us *privilegesService) UserRoleJoinAsRowData(ctx gctx.Context, upl []model.UserRoleJoin, allRoles []dbmodel.Roles) []basictable.RowData {
	roleOptions := make([]selectbox.SelectOptions, len(allRoles))
	for i, role := range allRoles {
		roleOptions[i].Label = role.Name
		roleOptions[i].Value = strconv.FormatInt(int64(role.ID), 10)
	}

	rows := make([]basictable.RowData, len(upl))
	for i, j := range upl {
		// selectAttrs := templ.Attributes{
		// 	"_": "on input trigger submit on closest <form/>",
		// }
		// selectBox := datainput.Select(fmt.Sprintf("%d-role-select", j.Users.ID), "role", strconv.FormatInt(int64(j.RoleID), 10), roleOptions, selectAttrs)
		formAttrs := templ.Attributes{
			"hx-put":     fmt.Sprintf("/roles/user/%d", j.Users.ID),
			"hx-trigger": "submit",
			"hx-swap":    "none",
		}
		form := form.Form(formAttrs)
		rows[i] = basictable.RowData{
			ID: "row-" + strconv.Itoa(i),
			Data: []basictable.CellData{
				{
					ID:    "n-" + strconv.Itoa(i),
					Width: 1,
					Body:  basiclabel.BasicLabel(fmt.Sprintf("%s %s", j.Users.FirstName, j.Users.LastName)),
				},
				{
					ID:    "em-" + strconv.Itoa(i),
					Width: 1,
					Body:  basiclabel.BasicLabel(j.Users.Email),
				},
				{
					ID:    "pr-" + strconv.Itoa(i),
					Width: 1,
					Body:  form,
				},
				{
					ID:    "ca-" + strconv.Itoa(i),
					Width: 1,
					Body:  basiclabel.BasicLabel(j.Users.CreatedAt.Format("2006-01-02")),
				},
			},
		}
	}
	return rows
}

func (us *privilegesService) JoinedRoleAsRowData(ctx gctx.Context, jpl []model.JoinedRole) []basictable.RowData {
	rows := make([]basictable.RowData, len(jpl))
	for i, j := range jpl {
		rows[i] = basictable.RowData{
			ID: "row-" + strconv.Itoa(i),
			Data: []basictable.CellData{
				{
					ID:    "n-" + strconv.Itoa(i),
					Width: 1,
					Body:  basiclabel.BasicLabel(j.RoleName),
				},
				{
					ID:    "pr-" + strconv.Itoa(i),
					Width: 1,
					Body:  basiclabel.BasicLabel(strconv.FormatInt(int64(j.RoleID), 10)),
				},
				{
					ID:    "ca-" + strconv.Itoa(i),
					Width: 1,
					Body:  basiclabel.BasicLabel(j.Privileges[0].CreatedAt.Format("2006-01-02")),
				},
			},
		}
	}
	return rows
}

func (us *privilegesService) JoinedPrivilegesAsRowData(ctx gctx.Context, jpl []model.JoinedPrivilegesRaw) []basictable.RowData {
	rows := make([]basictable.RowData, len(jpl))
	for i, p := range jpl {
		ca := p.PrivilegeCreatedAt
		caStr := "No Created At"
		if ca != nil {
			caStr = ca.String()
		}
		// xAttrs := templ.Attributes{
		// 	"class":      "fill-red-400 size-6 p-1 rounded-xs mx-auto cursor-pointer hover:bg-[#FFFFFF44]",
		// 	"hx-delete":  fmt.Sprintf("/roles-privileges/role/%d/privilege/%d", p.RoleID, p.PrivilegeID),
		// 	"hx-trigger": "click",
		// 	"hx-swap":    "none",
		// 	"_":          "on htmx:beforeRequest remove closest <tr/>",
		// }
		rows[i].ID = "row-" + strconv.Itoa(i)
		rows[i].Data = []basictable.CellData{
			{
				ID:    "role-" + strconv.Itoa(i),
				Width: 1,
				Body:  basiclabel.BasicLabel(p.RoleName),
			},
			{
				ID:    "pr-" + strconv.Itoa(i),
				Width: 1,
				Body:  basiclabel.BasicLabel(p.PrivilegeName),
			},
			{
				ID:    "ca-" + strconv.Itoa(i),
				Width: 1,
				Body:  basiclabel.BasicLabel(caStr),
			},
			{
				ID:    "del-" + strconv.Itoa(i),
				Width: 1,
				Body:  icon.X(),
			},
		}
	}
	return rows
}
