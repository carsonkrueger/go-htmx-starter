package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/tools"
)

type usersService struct {
	interfaces.IServiceContext
}

func NewUsersService(ctx interfaces.IServiceContext) *usersService {
	return &usersService{
		ctx,
	}
}

func (us *usersService) Login(email string, password string) (*string, error) {
	dao := us.DM().UsersDAO()
	user, err := dao.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(user.Password, "$")
	hash := tools.HashPassword(password, parts[0])

	if user.Password != hash {
		return nil, errors.New("Invalid password")
	}

	token, _ := tools.GenerateSalt()
	fullToken := fmt.Sprintf("%s$%d", token, user.ID)
	err = dao.UpdateAuthToken(user.ID, token)
	if err != nil {
		return nil, err
	}

	return &fullToken, nil
}
