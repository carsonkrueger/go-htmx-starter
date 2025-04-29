package services

import (
	"errors"
	"slices"
	"strconv"
	"sync"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/carsonkrueger/main/tools"
	"go.uber.org/zap"
)

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
