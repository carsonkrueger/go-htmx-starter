package auth

import (
	"database/sql"
	"time"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/carsonkrueger/main/models/auth"
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

func (dao *privilegesDAO) Insert(row *model.Privileges) (int64, error) {
	var res model.Privileges
	err := table.Privileges.INSERT(table.Privileges.Name).
		VALUES(postgres.String(row.Name)).
		RETURNING(table.Privileges.ID).
		Query(dao.db, &res)
	return res.ID, err
}

// Returns ID int64 if inserted.
// Parameter cols_update are the columns to be updated on conflict - if not provided, a few columns are updated
func (dao *privilegesDAO) Upsert(row *model.Privileges, colsUpdate ...postgres.ColumnAssigment) (int64, error) {
	if len(colsUpdate) == 0 {
		colsUpdate = []postgres.ColumnAssigment{
			table.Privileges.Name.SET(postgres.String(row.Name)),
			table.Privileges.UpdatedAt.SET(postgres.TimestampT(time.Now())),
		}
	}

	res, err := table.Privileges.
		INSERT(table.Privileges.Name).
		VALUES(row.Name).
		ON_CONFLICT(table.Privileges.Name).
		DO_UPDATE(postgres.SET(colsUpdate...)).
		RETURNING(table.Privileges.ID).
		Exec(dao.db)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	// if err then row was not inserted
	if err != nil {
		id = -1
	}

	return id, nil
}

func (dao *privilegesDAO) Update(row *model.Privileges) error {
	_, err := table.Privileges.
		UPDATE(table.Privileges.EXCLUDED.ID).
		MODEL(row).
		WHERE(table.Privileges.ID.EQ(postgres.Int(row.ID))).
		SET(table.Privileges.UpdatedAt.SET(postgres.TimestampT(time.Now()))).
		Exec(dao.db)
	if err != nil {
		return err
	}
	return nil
}

func (dao *privilegesDAO) Delete(id int64) error {
	_, err := table.Privileges.DELETE().WHERE(table.Privileges.ID.EQ(postgres.Int(id))).Exec(dao.db)
	if err != nil {
		return err
	}
	return nil
}

func (dao *privilegesDAO) GetAll() ([]*model.Privileges, error) {
	var rows []*model.Privileges
	err := table.Privileges.
		SELECT(table.Privileges.AllColumns).
		ORDER_BY(table.Privileges.ID.DESC()).
		Query(dao.db, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (dao *privilegesDAO) GetAllJoined() (map[int64][]model.Privileges, error) {
	var res []auth.JoinedPrivilegesRaw

	err := table.PrivilegeLevels.
		SELECT(
			table.PrivilegeLevels.ID.AS("JoinedPrivilegesRaw.LevelID"),
			table.PrivilegeLevels.Name.AS("JoinedPrivilegesRaw.LevelName"),
			table.Privileges.AllColumns,
		).
		FROM(table.PrivilegeLevels, table.Privileges).
		Query(dao.db, &res)
	if err != nil {
		return nil, err
	}

	mapping := make(map[int64][]model.Privileges)
	for _, row := range res {
		if _, exists := mapping[row.LevelID]; !exists {
			mapping[row.LevelID] = []model.Privileges{}
		}
		mapping[row.LevelID] = append(mapping[row.LevelID], row.Privileges)
	}

	return mapping, nil
}
