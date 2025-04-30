package services

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"sync"

	"github.com/a-h/templ"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/carsonkrueger/main/templates/datadisplay"
	"github.com/carsonkrueger/main/templates/datainput"
	"github.com/carsonkrueger/main/templates/partials"
	"github.com/carsonkrueger/main/tools"
	"go.uber.org/zap"
)

type PrivilegesService interface {
	BuildCache() error
	AddPermission(levelID int64, perms ...model.Privileges) error
	RemovePermission(levelID int64, perms ...model.Privileges) error
	CreateLevel(name string) error
	GetPermissions(levelID int64) []model.Privileges
	HasPermissionByID(levelID int64, permissionID int64) bool
	HasPermissionByName(levelID int64, permissionName string) bool
	SetUserPrivilegeLevel(levelID int64, userID int64) error
	UserPrivilegeLevelJoinAsRowData(upl []authModels.UserPrivilegeLevelJoin, allLevels []*model.PrivilegeLevels) []datadisplay.RowData
	JoinedPrivilegeLevelAsRowData(jpl []authModels.JoinedPrivilegeLevel) []datadisplay.RowData
	JoinedPrivilegesAsRowData(jpl []authModels.JoinedPrivilegesRaw) []datadisplay.RowData
}

type privilegesService struct {
	ServiceContext
	*permissionCache
}

func NewPrivilegesService(ctx ServiceContext, cache *permissionCache) *privilegesService {
	lgr := ctx.Lgr("NewPrivilegesService")
	if cache == nil {
		lgr.Panic("Permission cache is nil")
	}
	return &privilegesService{
		ctx,
		cache,
	}
}

type permissionCache struct {
	sync.RWMutex
	cache          authModels.PermissionCache
	levelNameCache authModels.LevelNameCache
}

func newPermissionCache() *permissionCache {
	return &permissionCache{
		cache: authModels.PermissionCache{},
	}
}

func (s *privilegesService) BuildCache() error {
	lgr := s.Lgr("BuildCache")
	lgr.Info("Called")
	dao := s.DM().PrivilegeDAO()

	joinedPrivileges, err := dao.GetAllJoined()
	if err != nil {
		lgr.Error(err.Error())
		return err
	}

	for _, jp := range joinedPrivileges {
		if jp.PrivilegeName == "" {
			continue
		}
		priv := model.Privileges{
			ID:   jp.PrivilegeID,
			Name: jp.PrivilegeName,
		}
		s.AddPermission(jp.LevelID, priv)
	}

	return nil
}

func (ps *privilegesService) AddPermission(levelID int64, perms ...model.Privileges) error {
	lgr := ps.Lgr("AddPermission")
	lgr.Info("Called")
	if len(perms) == 0 {
		return nil
	}

	tx, err := ps.DB().Begin()
	if err != nil {
		lgr.Error("Failed to begin transaction", zap.Error(err))
		return errors.New("Failed to begine transaction")
	}

	refd := tools.PtrSlice(perms)
	privDAO := ps.DM().PrivilegeDAO()
	if err = privDAO.UpsertMany(&refd, tx); err != nil {
		lgr.Error("Failed to insert privileges", zap.Error(err))
		return err
	}

	var joinRows []model.PrivilegeLevelsPrivileges
	for _, p := range perms {
		joinRows = append(joinRows, model.PrivilegeLevelsPrivileges{
			PrivilegeLevelID: levelID,
			PrivilegeID:      p.ID,
		})
	}
	levelsPrivsDAO := ps.DM().PrivilegeLevelsPrivilegesDAO()

	refdJoinRows := tools.PtrSlice(joinRows)
	if err = levelsPrivsDAO.UpsertMany(&refdJoinRows, tx); err != nil {
		lgr.Error("Failed to insert privilege level privileges", zap.Error(err), zap.Any("privilegeLevelPrivileges", refdJoinRows))
		return err
	}

	ps.RLock()
	permissions := ps.cache[levelID]
	ps.RUnlock()
	for _, p := range perms {
		if ps.HasPermissionByID(levelID, p.ID) {
			continue
		}
		ps.Lock()
		ps.cache[levelID] = append(permissions, p)
		ps.Unlock()
	}

	if err = tx.Commit(); err != nil {
		lgr.Error("Failed to commit transaction", zap.Error(err))
		return err
	}

	names := make([]string, len(perms))
	for i, p := range perms {
		names[i] = string(p.Name)
	}
	lgr.Info("Level:Privilege", zap.Strings(strconv.FormatInt(levelID, 10), names))

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

func (ps *privilegesService) GetPermissions(levelID int64) []model.Privileges {
	ps.RLock()
	defer ps.RUnlock()
	return ps.cache[levelID]
}

func (ps *privilegesService) HasPermissionByID(levelID int64, permissionID int64) bool {
	ps.RLock()
	permissions := ps.cache[levelID]
	ps.RUnlock()
	for _, p := range permissions {
		if p.ID == permissionID {
			return true
		}
	}
	return false
}

func (ps *privilegesService) HasPermissionByName(levelID int64, permissionName string) bool {
	ps.RLock()
	permissions := ps.cache[levelID]
	ps.RUnlock()
	for _, p := range permissions {
		if p.Name == permissionName {
			return true
		}
	}
	return false
}

func (ps *privilegesService) RemovePermission(levelID int64, privs ...model.Privileges) error {
	lgr := ps.Lgr("RemovePermission")
	lgr.Info("Called")

	ps.Lock()
	defer ps.Unlock()

	plpDAO := ps.DM().PrivilegeLevelsPrivilegesDAO()

	privileges := ps.cache[levelID]
	for _, removePerm := range privs {
		for i, perm := range privileges {
			if perm.ID == removePerm.ID {
				privileges = slices.Delete(privileges, i, i+1)
				pk := authModels.PrivilegeLevelsPrivilegesPrimaryKey{
					PrivilegeLevelID: levelID,
					PrivilegeID:      removePerm.ID,
				}
				if err := plpDAO.Delete(pk, ps.DB()); err != nil {
					lgr.Error("Failed to delete privilege", zap.Error(err))
				} else {
					lgr.Info("Removed privilege", zap.Int64("level", levelID), zap.Int64("privilege", removePerm.ID))
				}
				break
			}
		}
	}
	ps.cache[levelID] = privileges

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
