package interfaces

type IServiceManager interface {
	UsersService() IUsersService
	PrivilegesService() IPrivilegesService
}

type IUsersService interface {
	Login(email string, password string) (*string, error)
}

type IPrivilegesService interface {
	BuildCache() error
}
