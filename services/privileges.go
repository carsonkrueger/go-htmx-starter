package services

import (
	"fmt"
	"sync"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/models"
	"github.com/carsonkrueger/main/models/authModels"
)

type privilegesService struct {
	interfaces.IServiceContext
}

func NewPrivilegesService(ctx interfaces.IServiceContext) *privilegesService {
	return &privilegesService{
		ctx,
	}
}

type permissionCache struct {
	sync.RWMutex
	cache authModels.PermissionCache
}

func NewPermissionCache() *permissionCache {
	// pc.RWMutex = sync.RWMutex{}
	return &permissionCache{
		cache: authModels.PermissionCache{},
	}
}

// This function does not lock cache and returns a pointer which is unsafe
func (pc *permissionCache) unsafeGetPermissions(levelID int64) *[]model.Privileges {
	for _, item := range pc.cache {
		if item.A == levelID {
			return &item.B
		}
	}
	return nil
}

func (pc *permissionCache) AddPermission(levelID int64, perms ...model.Privileges) {
	pc.Lock()
	defer pc.Unlock()
	permissions := pc.unsafeGetPermissions(levelID)
	if permissions != nil {
		appended := append(*permissions, perms...)
		permissions = &appended
	} else {
		pc.cache = append(pc.cache, models.Pair[int64, []model.Privileges]{
			A: levelID,
			B: perms,
		})
	}
}

func (pc *permissionCache) GetPermissions(levelID int64) []model.Privileges {
	pc.RLock()
	defer pc.RUnlock()
	permissions := pc.unsafeGetPermissions(levelID)
	if permissions != nil {
		return *permissions
	}
	return []model.Privileges{}
}

func (pc *permissionCache) HasPermissionByID(levelID int64, permissionID int64) bool {
	pc.RLock()
	defer pc.RUnlock()
	permissions := pc.unsafeGetPermissions(levelID)
	for _, p := range *permissions {
		if p.ID == permissionID {
			return true
		}
	}
	return false
}

func (pc *permissionCache) HasPermissionByName(levelID int64, permissionName string) bool {
	pc.RLock()
	defer pc.RUnlock()
	permissions := pc.unsafeGetPermissions(levelID)
	fmt.Printf("%+v", permissions)
	for _, p := range *permissions {
		if p.Name == permissionName {
			return true
		}
	}
	return false
}

func (pc *permissionCache) SetPermissions(cache authModels.PermissionCache) {
	pc.Lock()
	defer pc.Unlock()
	pc.cache = cache
}

func (s *privilegesService) BuildCache() error {
	lgr := s.Lgr()
	cache := s.PC()
	dao := s.DM().PrivilegeDAO()

	joinedPrivileges, err := dao.GetAllJoined()
	lgr.Info(fmt.Sprintf("Cache: %+v", joinedPrivileges))
	if err != nil {
		return err
	}

	for _, jp := range joinedPrivileges {
		if jp.Privileges.Name == "" {
			continue
		}
		cache.AddPermission(jp.LevelID, jp.Privileges)
	}
	lgr.Info(fmt.Sprintf("Cache: %+v", cache))

	return nil
}
