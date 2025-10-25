package seeders

import (
	"database/sql"
	"fmt"

	"github.com/carsonkrueger/main/internal/constant"
	"github.com/carsonkrueger/main/internal/gen/go_starter_db/auth/table"
	"github.com/carsonkrueger/main/pkg/model/db/auth"
	"github.com/go-jet/jet/v2/postgres"
)

func SeedPermissions(db *sql.DB) error {
	allRoleNames := []string{
		constant.ADMIN_ROLE_NAME,
		constant.BASIC_ROLE_NAME,
	}
	adminRoleName := constant.ADMIN_ROLE_NAME

	newPrivileges := make([]auth.Roles, len(allRoleNames))
	for i, name := range allRoleNames {
		newPrivileges[i] = auth.Roles{Name: name}
	}

	fmt.Println("new privs", newPrivileges)

	_, err := table.Roles.
		INSERT(table.Roles.Name).
		MODELS(newPrivileges).
		ON_CONFLICT().
		DO_NOTHING().
		Exec(db)
	if err != nil {
		return err
	}

	var allPrivileges []auth.Privileges
	if err = table.Privileges.SELECT(table.Privileges.AllColumns).Query(db, &allPrivileges); err != nil {
		return err
	}

	if len(allPrivileges) == 0 {
		fmt.Println("No privileges found to seed, try running make live first.")
		return nil
	}

	var adminRole auth.Roles
	err = table.Roles.
		SELECT(table.Roles.AllColumns).
		WHERE(table.Roles.Name.EQ(postgres.String(adminRoleName))).
		Query(db, &adminRole)
	if err != nil {
		return err
	}

	rolesPrivileges := make([]auth.RolesPrivileges, len(allPrivileges))
	for i, p := range allPrivileges {
		rolesPrivileges[i] = auth.RolesPrivileges{
			RoleID:      adminRole.ID,
			PrivilegeID: p.ID,
		}
	}

	fmt.Println("Seeding privileges for admin role", rolesPrivileges)
	if len(rolesPrivileges) == 0 {
		fmt.Println("No roles-privileges found to seed, try running make live first.")
		return nil
	}

	_, err = table.RolesPrivileges.
		INSERT(table.RolesPrivileges.RoleID, table.RolesPrivileges.PrivilegeID).
		MODELS(rolesPrivileges).
		ON_CONFLICT().
		DO_NOTHING().
		Exec(db)
	if err != nil {
		return err
	}

	return nil
}

func UndoPermissions(db *sql.DB) error {
	_, err := table.Sessions.DELETE().
		WHERE(postgres.Bool(true)).
		Exec(db)
	if err != nil {
		return err
	}

	_, err = table.Users.DELETE().
		WHERE(postgres.Bool(true)).
		Exec(db)
	if err != nil {
		return err
	}

	_, err = table.RolesPrivileges.DELETE().
		WHERE(postgres.Bool(true)).
		Exec(db)
	if err != nil {
		return err
	}

	_, err = table.Privileges.DELETE().
		WHERE(postgres.Bool(true)).
		Exec(db)
	if err != nil {
		return err
	}

	_, err = table.Roles.DELETE().
		WHERE(postgres.Bool(true)).
		Exec(db)
	if err != nil {
		return err
	}

	return nil
}
