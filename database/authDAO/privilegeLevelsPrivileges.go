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

type privilegeLevelsPrivilegessDAO struct {
	db *sql.DB
}

func NewPrivilegeLevelsPrivilegesDAO(db *sql.DB) *privilegeLevelsPrivilegessDAO {
	return &privilegeLevelsPrivilegessDAO{
		db,
	}
}

func (dao *privilegeLevelsPrivilegessDAO) GetById(id authModels.PrivilegeLevelsPrivilegesPrimaryKey) (*model.PrivilegeLevelsPrivileges, error) {
	var row model.PrivilegeLevelsPrivileges
	err := table.PrivilegeLevelsPrivileges.
		SELECT(table.PrivilegeLevelsPrivileges.AllColumns).
		FROM(table.PrivilegeLevelsPrivileges).
		WHERE(table.PrivilegeLevelsPrivileges.PrivilegeID.EQ(postgres.Int(id.PrivilegeID)).
			AND(table.PrivilegeLevelsPrivileges.PrivilegeLevelID.EQ(postgres.Int(id.PrivilegeLevelID)))).
		LIMIT(1).
		Query(dao.db, &row)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (dao *privilegeLevelsPrivilegessDAO) Insert(row *model.PrivilegeLevelsPrivileges) error {
	return table.PrivilegeLevelsPrivileges.
		INSERT(table.PrivilegeLevelsPrivileges.AllColumns).
		MODEL(row).
		Query(dao.db, row)
}

func (dao *privilegeLevelsPrivilegessDAO) InsertMany(rows []*model.PrivilegeLevelsPrivileges) error {
	if len(rows) == 0 {
		return nil
	}
	return table.PrivilegeLevelsPrivileges.
		INSERT(table.PrivilegeLevelsPrivileges.AllColumns).
		MODELS(rows).
		Query(dao.db, &rows)
}

// Returns ID int64 if inserted.
// Parameter cols_update are the columns to be updated on conflict - if not provided, a few columns are updated
func (dao *privilegeLevelsPrivilegessDAO) Upsert(row *model.PrivilegeLevelsPrivileges, colsUpdate ...postgres.ColumnAssigment) error {
	if len(colsUpdate) == 0 {
		colsUpdate = []postgres.ColumnAssigment{
			table.PrivilegeLevelsPrivileges.PrivilegeID.SET(table.PrivilegeLevelsPrivileges.PrivilegeID),
			table.PrivilegeLevelsPrivileges.PrivilegeLevelID.SET(table.PrivilegeLevelsPrivileges.PrivilegeLevelID),
		}
	}

	return table.PrivilegeLevelsPrivileges.
		INSERT(table.PrivilegeLevelsPrivileges.AllColumns.Except(table.PrivilegeLevelsPrivileges.CreatedAt, table.PrivilegeLevelsPrivileges.UpdatedAt)).
		MODEL(row).
		ON_CONFLICT(table.PrivilegeLevelsPrivileges.PrivilegeID, table.PrivilegeLevelsPrivileges.PrivilegeLevelID).
		DO_UPDATE(postgres.SET(colsUpdate...)).
		RETURNING(table.PrivilegeLevelsPrivileges.PrivilegeID, table.PrivilegeLevelsPrivileges.PrivilegeLevelID).
		Query(dao.db, row)
}

func (dao *privilegeLevelsPrivilegessDAO) UpsertMany(rows []*model.PrivilegeLevelsPrivileges, colsUpdate ...postgres.ColumnAssigment) error {
	if len(colsUpdate) == 0 {
		colsUpdate = []postgres.ColumnAssigment{
			table.PrivilegeLevelsPrivileges.PrivilegeID.SET(table.PrivilegeLevelsPrivileges.PrivilegeID),
			table.PrivilegeLevelsPrivileges.PrivilegeLevelID.SET(table.PrivilegeLevelsPrivileges.PrivilegeLevelID),
		}
	}

	now := time.Now()
	for _, r := range rows {
		r.UpdatedAt = &now
	}

	return table.PrivilegeLevelsPrivileges.
		INSERT(table.PrivilegeLevelsPrivileges.EXCLUDED.AllColumns.Except(table.PrivilegeLevelsPrivileges.CreatedAt, table.PrivilegeLevelsPrivileges.UpdatedAt)).
		MODELS(rows).
		ON_CONFLICT(table.PrivilegeLevelsPrivileges.PrivilegeID, table.PrivilegeLevelsPrivileges.PrivilegeLevelID).
		DO_UPDATE(postgres.SET(colsUpdate...)).
		RETURNING(table.PrivilegeLevelsPrivileges.PrivilegeID, table.PrivilegeLevelsPrivileges.PrivilegeLevelID).
		Query(dao.db, &rows)
}

func (dao *privilegeLevelsPrivilegessDAO) Update(row *model.PrivilegeLevelsPrivileges) error {
	row.UpdatedAt = tools.Ptr(time.Now())
	_, err := table.PrivilegeLevelsPrivileges.
		UPDATE(table.PrivilegeLevelsPrivileges.AllColumns.Except(table.PrivilegeLevelsPrivileges.CreatedAt, table.PrivilegeLevelsPrivileges.UpdatedAt)).
		MODEL(row).
		WHERE(table.PrivilegeLevelsPrivileges.PrivilegeID.EQ(postgres.Int(row.PrivilegeID)).
			AND(table.PrivilegeLevelsPrivileges.PrivilegeLevelID.EQ(postgres.Int(row.PrivilegeLevelID)))).
		Exec(dao.db)
	return err
}

func (dao *privilegeLevelsPrivilegessDAO) Delete(id authModels.PrivilegeLevelsPrivilegesPrimaryKey) error {
	_, err := table.PrivilegeLevelsPrivileges.DELETE().
		WHERE(table.PrivilegeLevelsPrivileges.PrivilegeID.EQ(postgres.Int(id.PrivilegeID)).
			AND(table.PrivilegeLevelsPrivileges.PrivilegeLevelID.EQ(postgres.Int(id.PrivilegeLevelID)))).
		Exec(dao.db)
	if err != nil {
		return err
	}
	return nil
}

func (dao *privilegeLevelsPrivilegessDAO) GetAll() (*[]model.PrivilegeLevelsPrivileges, error) {
	var rows []model.PrivilegeLevelsPrivileges
	err := table.PrivilegeLevelsPrivileges.
		SELECT(table.PrivilegeLevelsPrivileges.AllColumns).
		ORDER_BY(table.PrivilegeLevelsPrivileges.PrivilegeLevelID.DESC()).
		Query(dao.db, &rows)
	if err != nil {
		return nil, err
	}
	return &rows, nil
}
