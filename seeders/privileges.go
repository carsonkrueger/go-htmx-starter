package seeders

import (
	"database/sql"

	"github.com/carsonkrueger/main/constant"
	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/go-jet/jet/v2/postgres"
)

func SeedPermissions(db *sql.DB) error {
	allPrivilegeLevelNames := []string{
		constant.ADMIN_LEVEL_NAME,
		constant.BASIC_LEVEL_NAME,
	}
	adminLevelName := constant.ADMIN_LEVEL_NAME

	newPrivileges := make([]model.PrivilegeLevels, len(allPrivilegeLevelNames))
	for i, name := range allPrivilegeLevelNames {
		newPrivileges[i] = model.PrivilegeLevels{Name: name}
	}

	_, err := table.PrivilegeLevels.
		INSERT(table.PrivilegeLevels.Name).
		MODELS(newPrivileges).
		ON_CONFLICT().
		DO_NOTHING().
		Exec(db)
	if err != nil {
		return err
	}

	var allLevels []model.PrivilegeLevels
	if err = table.PrivilegeLevels.SELECT(table.PrivilegeLevels.AllColumns).Query(db, &allLevels); err != nil {
		return err
	}

	var allPrivileges []model.Privileges
	if err = table.Privileges.SELECT(table.Privileges.AllColumns).Query(db, &allPrivileges); err != nil {
		return err
	}

	var adminLevel model.PrivilegeLevels
	err = table.PrivilegeLevels.
		SELECT(table.PrivilegeLevels.AllColumns).
		WHERE(table.PrivilegeLevels.Name.EQ(postgres.String(adminLevelName))).
		Query(db, &adminLevel)
	if err != nil {
		return err
	}

	privilegeLevelPrivileges := make([]model.PrivilegeLevelsPrivileges, len(allPrivileges))
	for i, p := range allPrivileges {
		privilegeLevelPrivileges[i] = model.PrivilegeLevelsPrivileges{
			PrivilegeLevelID: adminLevel.ID,
			PrivilegeID:      p.ID,
		}
	}

	_, err = table.PrivilegeLevelsPrivileges.
		INSERT(table.PrivilegeLevelsPrivileges.PrivilegeLevelID, table.PrivilegeLevelsPrivileges.PrivilegeID).
		MODELS(privilegeLevelPrivileges).
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

	_, err = table.PrivilegeLevelsPrivileges.DELETE().
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

	_, err = table.PrivilegeLevels.DELETE().
		WHERE(postgres.Bool(true)).
		Exec(db)
	if err != nil {
		return err
	}

	return nil
}
