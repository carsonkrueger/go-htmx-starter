package services

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/carsonkrueger/main/constant"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/carsonkrueger/main/tools"
	"go.uber.org/zap"
)

type UsersService interface {
	Login(email string, password string, req *http.Request) (*string, error)
	Logout(id int64, token string) error
	LogoutRequest(req *http.Request) error
	GetAuthParts(req *http.Request) (string, int64, error)
}

type usersService struct {
	ServiceContext
}

func NewUsersService(ctx ServiceContext) *usersService {
	return &usersService{
		ctx,
	}
}

func (us *usersService) Login(email string, password string, req *http.Request) (*string, error) {
	lgr := us.Lgr("Login")
	lgr.Info("Called")
	dao := us.DM().UsersDAO()

	user, err := dao.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	go us.LogoutRequest(req)

	parts := strings.Split(user.Password, "$")
	hash := tools.HashPassword(password, parts[0])

	if user.Password != hash {
		return nil, errors.New("Invalid password")
	}

	token, _ := tools.GenerateToken(32)
	fullToken := fmt.Sprintf("%s$%d", token, user.ID)

	row := &model.Sessions{
		UserID: user.ID,
		Token:  token,
	}
	sesDAO := us.DM().SessionsDAO()
	if err = sesDAO.Insert(row, us.DB()); err != nil {
		return nil, err
	}

	return &fullToken, nil
}

func (us *usersService) Logout(id int64, token string) error {
	lgr := us.Lgr("Logout")
	lgr.Info("Logging out", zap.Int64("user id", id))

	key := authModels.SessionsPrimaryKey{
		UserID:    id,
		AuthToken: token,
	}
	sesDAO := us.DM().SessionsDAO()
	return sesDAO.Delete(key, us.DB())
}

func (us *usersService) LogoutRequest(req *http.Request) error {
	lgr := us.Lgr("LogoutRequest")
	lgr.Info("Called")

	if req == nil {
		return errors.New("missing request")
	}
	token, id, err := us.GetAuthParts(req)

	if err != nil {
		return err
	}
	return us.Logout(id, token)
}

func (us *usersService) GetAuthParts(req *http.Request) (string, int64, error) {
	// lgr := us.Lgr("GetAuthParts")
	// lgr.Info("Called")

	cookie, err := tools.GetAuthCookie(req)
	if err != nil {
		return "", 0, err
	}

	token, id, err := tools.GetAuthParts(cookie)
	if err != nil {
		req.Header.Del(constant.AUTH_TOKEN_KEY)
		return "", 0, err
	}

	return token, id, nil
}
