package services

import (
	"github.com/carsonkrueger/main/interfaces"
)

type privilegesService struct {
	interfaces.IServiceContext
}

func NewPrivilegesService(ctx interfaces.IServiceContext) *privilegesService {
	return &privilegesService{
		ctx,
	}
}

func (s *privilegesService) BuildCache() error {
	s.Lgr().Info("BuildCache() called")
	cache := s.PC()
	dao := s.DM().PrivilegeDAO()
	mapping, err := dao.GetAllJoined()
	if err != nil {
		return err
	}
	// s.Lgr().Info(fmt.Sprintf("Joined Permissions Array: %v", rows))
	cache.SetPermissions(mapping)
	return nil
}
