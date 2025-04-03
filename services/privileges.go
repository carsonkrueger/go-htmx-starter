package services

import (
	"strconv"
	"sync"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/carsonkrueger/main/tools"
	"go.uber.org/zap"
)

type privilegesService struct {
	interfaces.IServiceContext
	*permissionCache
}

func NewPrivilegesService(ctx interfaces.IServiceContext, cache *permissionCache) *privilegesService {
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
	cache authModels.PermissionCache
}

func NewPermissionCache() *permissionCache {
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
		return err
	}

	for _, jp := range *joinedPrivileges {
		if jp.Privileges.Name == "" {
			continue
		}
		s.AddPermission(jp.LevelID, jp.Privileges)
	}

	return nil
}

func (ps *privilegesService) AddPermission(levelID int64, perms ...model.Privileges) {
	lgr := ps.Lgr("AddPermission")
	if len(perms) == 0 {
		return
	}

	refd := tools.PtrSlice(perms)
	err := ps.DM().PrivilegeDAO().UpsertMany(refd)
	if err != nil {
		lgr.Error("Failed to insert privileges", zap.Error(err))
		return
	}

	ps.Lock()
	defer ps.Unlock()

	permissions := ps.cache[levelID]
	ps.cache[levelID] = append(permissions, perms...)

	names := make([]string, len(perms))
	for i, p := range perms {
		names[i] = string(p.Name)
	}
	lgr.Info("Level:Privilege", zap.Strings(strconv.FormatInt(levelID, 10), names))
}

func (ps *privilegesService) GetPermissions(levelID int64) []model.Privileges {
	ps.RLock()
	defer ps.RUnlock()
	return ps.cache[levelID]
}

func (ps *privilegesService) HasPermissionByID(levelID int64, permissionID int64) bool {
	ps.RLock()
	defer ps.RUnlock()
	permissions := ps.cache[levelID]
	for _, p := range permissions {
		if p.ID == permissionID {
			return true
		}
	}
	return false
}

func (ps *privilegesService) HasPermissionByName(levelID int64, permissionName string) bool {
	ps.RLock()
	defer ps.RUnlock()
	permissions := ps.cache[levelID]
	for _, p := range permissions {
		if p.Name == permissionName {
			return true
		}
	}
	return false
}
