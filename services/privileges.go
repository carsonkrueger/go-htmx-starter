package services

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"github.com/carsonkrueger/main/templates/datainput"
	"github.com/carsonkrueger/main/templates/partials"
	"go.uber.org/zap"
)

type PrivilegesService interface {
	CreatePrivilegeAssociation(levelID int64, privID int64) error
	DeletePrivilegeAssociation(levelID int64, privID int64) error
	CreateLevel(name string) error
	HasPermissionByID(levelID int64, permissionID int64) bool
	SetUserPrivilegeLevel(levelID int64, userID int64) error
	UserPrivilegeLevelJoinAsRowData(upl []authModels.UserPrivilegeLevelJoin, allLevels []*model.PrivilegeLevels) []datadisplay.RowData
	JoinedPrivilegeLevelAsRowData(jpl []authModels.JoinedPrivilegeLevel) []datadisplay.RowData
	JoinedPrivilegesAsRowData(jpl []authModels.JoinedPrivilegesRaw) []datadisplay.RowData
}

type privilegesService struct {
	ServiceContext
}

func NewPrivilegesService(ctx ServiceContext) *privilegesService {
	return &privilegesService{ctx}
}

func (ps *privilegesService) CreatePrivilegeAssociation(levelID int64, privID int64) error {
	lgr := ps.Lgr("AddPermission")
	lgr.Info("Level:Privilege", zap.String(strconv.FormatInt(levelID, 10), strconv.FormatInt(privID, 10)))

	joinRow := model.PrivilegeLevelsPrivileges{
		PrivilegeLevelID: levelID,
		PrivilegeID:      privID,
	}
	levelsPrivsDAO := ps.DM().PrivilegeLevelsPrivilegesDAO()

	if err := levelsPrivsDAO.Upsert(&joinRow, ps.DB()); err != nil {
		lgr.Error("Failed to insert privilege level privileges", zap.Error(err), zap.Any("privilegeLevelPrivileges", joinRow))
		return err
	}

	return nil
}

func (ps *privilegesService) CreateLevel(name string) error {
	lgr := ps.Lgr("CreateLevel")
	lgr.Info("Called")

	row := model.PrivilegeLevels{
		Name: name,
	}
	levelsDAO := ps.DM().PrivilegeLevelsDAO()
	if err := levelsDAO.Insert(&row, ps.DB()); err != nil {
		lgr.Error("Failed to create level", zap.Error(err))
		return errors.New("Failed to create level")
	}
	return nil
}

func (ps *privilegesService) HasPermissionByID(levelID int64, permissionID int64) bool {
	pk := authModels.PrivilegeLevelsPrivilegesPrimaryKey{
		PrivilegeID:      permissionID,
		PrivilegeLevelID: levelID,
	}
	row, err := ps.DM().PrivilegeLevelsPrivilegesDAO().GetOne(pk, ps.DB())
	return row != nil && err == nil
}

func (ps *privilegesService) DeletePrivilegeAssociation(levelID int64, privID int64) error {
	lgr := ps.Lgr("DeletePrivilegeAssociation")
	lgr.Info("Called")

	pk := authModels.PrivilegeLevelsPrivilegesPrimaryKey{
		PrivilegeID:      privID,
		PrivilegeLevelID: levelID,
	}
	if err := ps.DM().PrivilegeLevelsPrivilegesDAO().Delete(pk, ps.DB()); err != nil {
		return nil
	}

	return nil
}

func (us *privilegesService) SetUserPrivilegeLevel(levelID int64, userID int64) error {
	lgr := us.Lgr("SetUserPrivilegeLevel")
	lgr.Info("Called")

	db := us.DB()
	levelDAO := us.DM().PrivilegeLevelsDAO()
	level, err := levelDAO.GetOne(levelID, db)
	if err != nil {
		lgr.Error("Error fetching privilege level", zap.Error(err))
		return err
	}

	userDAO := us.DM().UsersDAO()
	user, err := userDAO.GetOne(userID, db)
	if err != nil {
		lgr.Error("Error fetching user", zap.Error(err))
		return err
	}

	user.PrivilegeLevelID = level.ID
	if err := userDAO.Update(user, user.ID, db); err != nil {
		lgr.Error("Error updating user", zap.Error(err))
		return err
	}

	return nil
}

func (us *privilegesService) UserPrivilegeLevelJoinAsRowData(upl []authModels.UserPrivilegeLevelJoin, allLevels []*model.PrivilegeLevels) []datadisplay.RowData {
	levelOptions := make([]datainput.SelectOptions, len(allLevels))
	for i, level := range allLevels {
		levelOptions[i].Label = level.Name
		levelOptions[i].Value = strconv.FormatInt(level.ID, 10)
	}

	rows := make([]datadisplay.RowData, len(upl))
	for i, j := range upl {
		selectAttrs := templ.Attributes{
			"_": "on input trigger submit on closest <form/>",
		}
		selectBox := datainput.Select(fmt.Sprintf("%d-lvl-select", j.Users.ID), "privilege-level", strconv.FormatInt(j.PrivilegeLevelID, 10), levelOptions, selectAttrs)
		formAttrs := templ.Attributes{
			"hx-put":     fmt.Sprintf("/privilege-levels/user/%d", j.Users.ID),
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

func (us *privilegesService) JoinedPrivilegeLevelAsRowData(jpl []authModels.JoinedPrivilegeLevel) []datadisplay.RowData {
	rows := make([]datadisplay.RowData, len(jpl))
	for i, j := range jpl {
		rows[i] = datadisplay.RowData{
			ID: "row-" + strconv.Itoa(i),
			Data: []datadisplay.CellData{
				{
					ID:    "n-" + strconv.Itoa(i),
					Width: 1,
					Body:  datadisplay.Text(j.LevelName, datadisplay.MD),
				},
				{
					ID:    "pr-" + strconv.Itoa(i),
					Width: 1,
					Body:  datadisplay.Text(strconv.FormatInt(j.LevelID, 10), datadisplay.MD),
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

func (us *privilegesService) JoinedPrivilegesAsRowData(jpl []authModels.JoinedPrivilegesRaw) []datadisplay.RowData {
	rows := make([]datadisplay.RowData, len(jpl))
	for i, p := range jpl {
		ca := p.PrivilegeCreatedAt
		caStr := "No Created At"
		if ca != nil {
			caStr = ca.String()
		}
		xAttrs := templ.Attributes{
			"class":      "fill-red-400 size-6 p-1 rounded-xs mx-auto cursor-pointer hover:bg-[#FFFFFF44]",
			"hx-delete":  fmt.Sprintf("/privilege-levels-privileges/level/%d/privilege/%d", p.LevelID, p.PrivilegeID),
			"hx-trigger": "click",
			"hx-swap":    "none",
			"_":          "on htmx:beforeRequest remove closest <tr/>",
		}
		rows[i].ID = "row-" + strconv.Itoa(i)
		rows[i].Data = []datadisplay.CellData{
			{
				ID:    "lvl-" + strconv.Itoa(i),
				Width: 1,
				Body:  datadisplay.Text(p.LevelName, datadisplay.SM),
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
