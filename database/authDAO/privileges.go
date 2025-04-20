package authDAO

import (
	"database/sql"
	"time"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/go-jet/jet/v2/postgres"
)

type privilegesDAO struct {
	db *sql.DB
}

func NewPrivilegesDAO(db *sql.DB) *privilegesDAO {
	return &privilegesDAO{
		db,
	}
}

func (dao *privilegesDAO) Table() interfaces.IPostgresTable {
	return table.Privileges
}

func (dao *privilegesDAO) InsertCols() postgres.ColumnList {
	return table.Privileges.AllColumns.Except(
		table.Privileges.ID,
		table.Privileges.CreatedAt,
		table.Privileges.UpdatedAt,
	)
}

func (dao *privilegesDAO) UpdateCols() postgres.ColumnList {
	return table.Privileges.AllColumns.Except(
		table.Privileges.ID,
		table.Privileges.CreatedAt,
	)
}

func (dao *privilegesDAO) AllCols() postgres.ColumnList {
	return table.Privileges.AllColumns
}

func (dao *privilegesDAO) OnConflictCols() postgres.ColumnList {
	return []postgres.Column{table.Privileges.Name}
}

func (dao *privilegesDAO) UpdateOnConflictCols() []postgres.ColumnAssigment {
	return []postgres.ColumnAssigment{
		table.Privileges.Name.SET(table.Privileges.EXCLUDED.Name),
	}
}

func (dao *privilegesDAO) PKMatch(pk int64) postgres.BoolExpression {
	return table.Privileges.ID.EQ(postgres.Int(pk))
}

func (dao *privilegesDAO) GetUpdatedAt(row *model.Privileges) *time.Time {
	return row.UpdatedAt
}

func (dao *privilegesDAO) GetAllJoined() (*[]authModels.JoinedPrivilegesRaw, error) {
	var res []authModels.JoinedPrivilegesRaw

	err := table.PrivilegeLevels.
		SELECT(
			table.PrivilegeLevels.ID.AS("JoinedPrivilegesRaw.LevelID"),
			table.PrivilegeLevels.Name.AS("JoinedPrivilegesRaw.LevelName"),
			table.Privileges.AllColumns,
		).
		FROM(
			table.PrivilegeLevels.
				LEFT_JOIN(table.PrivilegeLevelsPrivileges, table.PrivilegeLevelsPrivileges.PrivilegeLevelID.EQ(table.PrivilegeLevels.ID)).
				LEFT_JOIN(table.Privileges, table.Privileges.ID.EQ(table.PrivilegeLevelsPrivileges.PrivilegeID)),
		).
		Query(dao.db, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
