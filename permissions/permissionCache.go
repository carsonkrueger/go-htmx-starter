package permissions

import (
	"sync"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
)

type permissionCache struct {
	sync.RWMutex
	cache map[int64][]model.Privileges
}

func NewPermissionCache() *permissionCache {
	// pc.RWMutex = sync.RWMutex{}
	return &permissionCache{
		cache: make(map[int64][]model.Privileges),
	}
}

func (pc *permissionCache) AddPermission(levelID int64, perm model.Privileges) {
	pc.Lock()
	defer pc.Unlock()
	pc.cache[levelID] = append(pc.cache[levelID], perm)
}

func (pc *permissionCache) GetPermissions(levelID int64) []model.Privileges {
	pc.RLock()
	defer pc.RUnlock()
	return pc.cache[levelID]
}

func (pc *permissionCache) HasPermissionByID(levelID int64, permissionID int64) bool {
	pc.RLock()
	defer pc.RUnlock()
	for _, p := range pc.cache[levelID] {
		if p.ID == permissionID {
			return true
		}
	}
	return false
}

func (pc *permissionCache) HasPermissionByName(levelID int64, permissionName string) bool {
	pc.RLock()
	defer pc.RUnlock()
	for _, p := range pc.cache[levelID] {
		if p.Name == permissionName {
			return true
		}
	}
	return false
}
