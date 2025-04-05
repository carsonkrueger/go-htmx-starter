package authDAO

import (
	"database/sql"
	"time"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/carsonkrueger/main/models/authModels"
	"github.com/carsonkrueger/main/tools"
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

func (dao *privilegesDAO) GetById(id int64) (*model.Privileges, error) {
	var row model.Privileges
	err := table.Privileges.SELECT(table.Privileges.AllColumns).
		FROM(table.Privileges).
		WHERE(table.Privileges.ID.EQ(postgres.Int(id))).
		LIMIT(1).
		Query(dao.db, &row)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (dao *privilegesDAO) Insert(row *model.Privileges) error {
	var res model.Privileges
	err := table.Privileges.
		INSERT(table.Privileges.AllColumns.Except(table.Privileges.ID, table.Privileges.CreatedAt, table.Privileges.UpdatedAt)).
		VALUES(postgres.String(row.Name)).
		RETURNING(table.Privileges.ID).
		Query(dao.db, res)
	return err
}

func (dao *privilegesDAO) InsertMany(rows []*model.Privileges) error {
	if len(rows) == 0 {
		return nil
	}
	return table.Privileges.
		INSERT(table.Privileges.AllColumns.Except(table.Privileges.ID, table.Privileges.CreatedAt, table.Privileges.UpdatedAt)).
		MODELS(rows).
		RETURNING(table.Privileges.ID).
		Query(dao.db, &rows)
}

// Returns ID int64 if inserted.
// Parameter cols_update are the columns to be updated on conflict - if not provided, a few columns are updated
func (dao *privilegesDAO) Upsert(row *model.Privileges, colsUpdate ...postgres.ColumnAssigment) error {
	if len(colsUpdate) == 0 {
		colsUpdate = []postgres.ColumnAssigment{
			table.Privileges.Name.SET(postgres.String(row.Name)),
		}
	}

	row.UpdatedAt = tools.Ptr(time.Now())

	return table.Privileges.
		INSERT(table.Privileges.AllColumns.Except(table.Privileges.ID, table.Privileges.CreatedAt, table.Privileges.UpdatedAt)).
		VALUES(row.Name).
		ON_CONFLICT(table.Privileges.Name).
		DO_UPDATE(postgres.SET(colsUpdate...)).
		RETURNING(table.Privileges.ID).
		Query(dao.db, row)
}

func (dao *privilegesDAO) UpsertMany(rows []*model.Privileges, colsUpdate ...postgres.ColumnAssigment) error {
	if len(colsUpdate) == 0 {
		colsUpdate = []postgres.ColumnAssigment{
			table.Privileges.Name.SET(table.Privileges.Name),
		}
	}

	now := time.Now()
	for _, r := range rows {
		r.UpdatedAt = &now
	}

	return table.Privileges.
		INSERT(table.Privileges.AllColumns.Except(table.Privileges.ID, table.Privileges.CreatedAt, table.Privileges.UpdatedAt)).
		MODELS(rows).
		ON_CONFLICT(table.Privileges.Name).
		DO_UPDATE(postgres.SET(colsUpdate...)).
		RETURNING(table.Privileges.ID).
		Query(dao.db, &rows)
}

func (dao *privilegesDAO) Update(row *model.Privileges) error {
	row.UpdatedAt = tools.Ptr(time.Now())
	_, err := table.Privileges.
		UPDATE(table.Privileges.EXCLUDED.ID).
		MODEL(row).
		WHERE(table.Privileges.ID.EQ(postgres.Int(row.ID))).
		SET(table.Privileges.UpdatedAt.SET(postgres.TimestampT(time.Now()))).
		Exec(dao.db)
	return err
}

func (dao *privilegesDAO) Delete(id int64) error {
	_, err := table.Privileges.DELETE().WHERE(table.Privileges.ID.EQ(postgres.Int(id))).Exec(dao.db)
	if err != nil {
		return err
	}
	return nil
}

func (dao *privilegesDAO) GetAll() (*[]model.Privileges, error) {
	var rows []model.Privileges
	err := table.Privileges.
		SELECT(table.Privileges.AllColumns).
		ORDER_BY(table.Privileges.ID.DESC()).
		Query(dao.db, &rows)
	if err != nil {
		return nil, err
	}
	return &rows, nil
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
