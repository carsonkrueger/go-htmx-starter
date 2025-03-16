package interfaces

type IServiceManager interface {
	UsersService() IUsersService
}

type IUsersService interface {
	Login(email string, password string) (*string, error)
	IsPermitted(userId int64, permissionName string) bool
}
