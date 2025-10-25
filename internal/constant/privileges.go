package constant

type PrivilegeName string

// These should NEVER be changed
const (
	// Users
	UsersRead   PrivilegeName = "users:read"
	UsersCreate PrivilegeName = "users:create"
	UsersUpdate PrivilegeName = "users:update"
	UsersDelete PrivilegeName = "users:delete"
	UsersList   PrivilegeName = "users:list"

	// Roles
	RolesRead   PrivilegeName = "roles:read"
	RolesCreate PrivilegeName = "roles:create"
	RolesUpdate PrivilegeName = "roles:update"
	RolesDelete PrivilegeName = "roles:delete"
	RolesList   PrivilegeName = "roles:list"

	// Privileges
	PrivilegesRead   PrivilegeName = "privileges:read"
	PrivilegesCreate PrivilegeName = "privileges:create"
	PrivilegesUpdate PrivilegeName = "privileges:update"
	PrivilegesDelete PrivilegeName = "privileges:delete"
	PrivilegesList   PrivilegeName = "privileges:list"
)
