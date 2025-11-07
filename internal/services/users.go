package services

import (
	gctx "context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/carsonkrueger/main/internal/constant"
	"github.com/carsonkrueger/main/internal/context"
	dbmodel "github.com/carsonkrueger/main/pkg/db/auth/model"
	"github.com/carsonkrueger/main/pkg/model"
	"github.com/carsonkrueger/main/pkg/util"
)

type usersService struct {
	*context.AppContext
}

func NewUsersService(ctx *context.AppContext) *usersService {
	return &usersService{
		ctx,
	}
}

func (us *usersService) Login(ctx gctx.Context, email string, password string, req *http.Request) (*string, error) {
	dao := us.DM().UsersDAO()

	user, err := dao.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	go us.LogoutRequest(ctx, req)

	parts := strings.Split(user.Password, "$")
	hash := util.HashPassword(password, parts[0])

	if user.Password != hash {
		return nil, errors.New("Invalid password")
	}

	token, _ := util.GenerateToken(32)
	fullToken := fmt.Sprintf("%s$%d", token, user.ID)

	row := &dbmodel.Sessions{
		UserID: user.ID,
		Token:  token,
	}
	sesDAO := us.DM().SessionsDAO()
	if err = sesDAO.Insert(ctx, row); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &fullToken, nil
}

func (us *usersService) Logout(ctx gctx.Context, id int64, token string) error {
	key := model.SessionsPrimaryKey{
		UserID:    id,
		AuthToken: token,
	}
	sesDAO := us.DM().SessionsDAO()
	if err := sesDAO.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

func (us *usersService) LogoutRequest(ctx gctx.Context, req *http.Request) error {
	if req == nil {
		return errors.New("missing request")
	}
	token, id, err := us.GetAuthParts(ctx, req)

	if err != nil {
		return fmt.Errorf("failed to get auth parts: %w", err)
	}
	return us.Logout(ctx, id, token)
}

func (us *usersService) GetAuthParts(ctx gctx.Context, req *http.Request) (string, int64, error) {
	cookie, err := util.GetAuthCookie(req, constant.AUTH_TOKEN_KEY)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get auth cookie: %w", err)
	}

	token, id, err := util.GetAuthParts(cookie)
	if err != nil {
		req.Header.Del(constant.AUTH_TOKEN_KEY)
		return "", 0, fmt.Errorf("failed to get auth parts: %w", err)
	}

	return token, id, nil
}
