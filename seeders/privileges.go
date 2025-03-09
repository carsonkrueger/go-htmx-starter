package seeders

import (
	"database/sql"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/carsonkrueger/main/internal"
	"github.com/go-jet/jet/v2/postgres"
)

var seeds []model.Privileges = []model.Privileges{
	internal.HelloWorldGet,
	internal.HelloWorldGet2,
}

type SeedPrivilegeLevel struct {
	ID         int64
	Name       string
	Privileges []model.Privileges
}

var basicLevel SeedPrivilegeLevel = SeedPrivilegeLevel{
	ID:   1000,
	Name: "basic",
	Privileges: []model.Privileges{
		internal.HelloWorldGet,
	},
}

var adminLevel SeedPrivilegeLevel = SeedPrivilegeLevel{
	ID:         1001,
	Name:       "admin",
	Privileges: append(basicLevel.Privileges, internal.HelloWorldGet2),
}

func SeedPermissions(db *sql.DB) error {
	_, err := table.Privileges.INSERT(table.Privileges.ID, table.Privileges.Name).
		MODELS(seeds).
		ON_CONFLICT().
		DO_NOTHING().
		Exec(db)

	if err != nil {
		return err
	}

	levels := []SeedPrivilegeLevel{
		basicLevel,
		adminLevel,
	}

	for _, level := range levels {
		_, err := table.PrivilegeLevels.INSERT(table.PrivilegeLevels.ID, table.PrivilegeLevels.Name).
			MODEL(model.PrivilegeLevels{
				ID:   level.ID,
				Name: level.Name,
			}).
			ON_CONFLICT().
			DO_NOTHING().
			Exec(db)

		if err != nil {
			return err
		}

		for _, privilege := range level.Privileges {
			_, err := table.PrivilegeLevelsPrivileges.INSERT(table.PrivilegeLevelsPrivileges.PrivilegeLevelID, table.PrivilegeLevelsPrivileges.PrivilegeID).
				MODEL(model.PrivilegeLevelsPrivileges{
					PrivilegeLevelID: level.ID,
					PrivilegeID:      privilege.ID,
				}).
				ON_CONFLICT().
				DO_NOTHING().
				Exec(db)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func UndoPermissions(db *sql.DB) error {
	_, err := table.Users.DELETE().
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
