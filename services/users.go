package services

import (
	gctx "context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/carsonkrueger/main/constant"
	"github.com/carsonkrueger/main/context"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/models/auth_models"
	"github.com/carsonkrueger/main/tools"
	"go.uber.org/zap"
)

type usersService struct {
	context.AppContext
}

func NewUsersService(ctx context.AppContext) *usersService {
	return &usersService{
		ctx,
	}
}

func (us *usersService) Login(ctx gctx.Context, email string, password string, req *http.Request) (*string, error) {
	lgr := us.Lgr("Login")
	lgr.Info("Called")
	dao := us.DM().UsersDAO()

	user, err := dao.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	go us.LogoutRequest(ctx, req)

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
	if err = sesDAO.Insert(ctx, row); err != nil {
		return nil, err
	}

	return &fullToken, nil
}

func (us *usersService) Logout(ctx gctx.Context, id int64, token string) error {
	lgr := us.Lgr("Logout")
	lgr.Info("Logging out", zap.Int64("user id", id))

	key := auth_models.SessionsPrimaryKey{
		UserID:    id,
		AuthToken: token,
	}
	sesDAO := us.DM().SessionsDAO()
	return sesDAO.Delete(ctx, key)
}

func (us *usersService) LogoutRequest(ctx gctx.Context, req *http.Request) error {
	lgr := us.Lgr("LogoutRequest")
	lgr.Info("Called")

	if req == nil {
		return errors.New("missing request")
	}
	token, id, err := us.GetAuthParts(ctx, req)

	if err != nil {
		return err
	}
	return us.Logout(ctx, id, token)
}

func (us *usersService) GetAuthParts(ctx gctx.Context, req *http.Request) (string, int64, error) {
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
