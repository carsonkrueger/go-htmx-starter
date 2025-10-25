package services

import (
	gctx "context"
	"errors"
	"fmt"
	"strconv"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/internal/context"
	"github.com/carsonkrueger/main/internal/templates/datadisplay"
	"github.com/carsonkrueger/main/internal/templates/datainput"
	"github.com/carsonkrueger/main/internal/templates/partials"
	model1 "github.com/carsonkrueger/main/pkg/model"
	"github.com/carsonkrueger/main/pkg/model/db/auth"
	"go.uber.org/zap"
)

type privilegesService struct {
	context.AppContext
}

func NewPrivilegesService(ctx context.AppContext) *privilegesService {
	return &privilegesService{ctx}
}

func (ps *privilegesService) CreatePrivilegeAssociation(ctx gctx.Context, roleID int16, privID int64) error {
	lgr := ps.Lgr("AddPermission")
	lgr.Info("Role:Privilege", zap.String(strconv.FormatInt(int64(roleID), 10), strconv.FormatInt(int64(privID), 10)))

	joinRow := auth.RolesPrivileges{
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

	row := auth.Roles{
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
	pks := make([]model1.RolesPrivilegesPrimaryKey, len(privIDs))
	for i, privID := range privIDs {
		pks[i] = model1.RolesPrivilegesPrimaryKey{
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

	pk := model1.RolesPrivilegesPrimaryKey{
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

func (us *privilegesService) UserRoleJoinAsRowData(ctx gctx.Context, upl []model1.UserRoleJoin, allRoles []auth.Roles) []datadisplay.RowData {
	roleOptions := make([]datainput.SelectOptions, len(allRoles))
	for i, role := range allRoles {
		roleOptions[i].Label = role.Name
		roleOptions[i].Value = strconv.FormatInt(int64(role.ID), 10)
	}

	rows := make([]datadisplay.RowData, len(upl))
	for i, j := range upl {
		selectAttrs := templ.Attributes{
			"_": "on input trigger submit on closest <form/>",
		}
		selectBox := datainput.Select(fmt.Sprintf("%d-role-select", j.Users.ID), "role", strconv.FormatInt(int64(j.RoleID), 10), roleOptions, selectAttrs)
		formAttrs := templ.Attributes{
			"hx-put":     fmt.Sprintf("/roles/user/%d", j.Users.ID),
			"hx-trigger": "submit",
			"hx-swap":    "none",
		}
		form := partials.FormBasic(selectBox, formAttrs)
		rows[i] = datadisplay.RowData{
			ID: "row-" + strconv.Itoa(i),
			Data: []datadisplay.CellData{
				{
					ID:    "n-" + strconv.Itoa(i),
					Width: 1,
					Body:  datadisplay.Text(fmt.Sprintf("%s %s", j.FirstName, j.LastName), datadisplay.MD),
				},
				{
					ID:    "em-" + strconv.Itoa(i),
					Width: 1,
					Body:  datadisplay.Text(j.Email, datadisplay.MD),
				},
				{
					ID:    "pr-" + strconv.Itoa(i),
					Width: 1,
					Body:  form,
				},
				{
					ID:    "ca-" + strconv.Itoa(i),
					Width: 1,
					Body:  datadisplay.Text(j.CreatedAt.Format("2006-01-02"), datadisplay.MD),
				},
			},
		}
	}
	return rows
}

func (us *privilegesService) JoinedRoleAsRowData(ctx gctx.Context, jpl []model1.JoinedRole) []datadisplay.RowData {
	rows := make([]datadisplay.RowData, len(jpl))
	for i, j := range jpl {
		rows[i] = datadisplay.RowData{
			ID: "row-" + strconv.Itoa(i),
			Data: []datadisplay.CellData{
				{
					ID:    "n-" + strconv.Itoa(i),
					Width: 1,
					Body:  datadisplay.Text(j.RoleName, datadisplay.MD),
				},
				{
					ID:    "pr-" + strconv.Itoa(i),
					Width: 1,
					Body:  datadisplay.Text(strconv.FormatInt(int64(j.RoleID), 10), datadisplay.MD),
				},
				{
					ID:    "ca-" + strconv.Itoa(i),
					Width: 1,
					Body:  datadisplay.Text(j.Privileges[0].CreatedAt.Format("2006-01-02"), datadisplay.MD),
				},
			},
		}
	}
	return rows
}

func (us *privilegesService) JoinedPrivilegesAsRowData(ctx gctx.Context, jpl []model1.JoinedPrivilegesRaw) []datadisplay.RowData {
	rows := make([]datadisplay.RowData, len(jpl))
	for i, p := range jpl {
		ca := p.PrivilegeCreatedAt
		caStr := "No Created At"
		if ca != nil {
			caStr = ca.String()
		}
		xAttrs := templ.Attributes{
			"class":      "fill-red-400 size-6 p-1 rounded-xs mx-auto cursor-pointer hover:bg-[#FFFFFF44]",
			"hx-delete":  fmt.Sprintf("/roles-privileges/role/%d/privilege/%d", p.RoleID, p.PrivilegeID),
			"hx-trigger": "click",
			"hx-swap":    "none",
			"_":          "on htmx:beforeRequest remove closest <tr/>",
		}
		rows[i].ID = "row-" + strconv.Itoa(i)
		rows[i].Data = []datadisplay.CellData{
			{
				ID:    "role-" + strconv.Itoa(i),
				Width: 1,
				Body:  datadisplay.Text(p.RoleName, datadisplay.SM),
			},
			{
				ID:    "pr-" + strconv.Itoa(i),
				Width: 1,
				Body:  datadisplay.Text(p.PrivilegeName, datadisplay.SM),
			},
			{
				ID:    "ca-" + strconv.Itoa(i),
				Width: 1,
				Body:  datadisplay.Text(caStr, datadisplay.MD),
			},
			{
				ID:    "del-" + strconv.Itoa(i),
				Width: 1,
				Body:  datadisplay.X(xAttrs),
			},
		}
	}
	return rows
}
