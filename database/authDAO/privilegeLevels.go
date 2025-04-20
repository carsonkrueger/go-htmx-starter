package authDAO

import (
	"database/sql"
	"time"

	"github.com/carsonkrueger/main/gen/go_db/auth/model"
	"github.com/carsonkrueger/main/gen/go_db/auth/table"
	"github.com/carsonkrueger/main/interfaces"
	"github.com/carsonkrueger/main/tools"
	"github.com/go-jet/jet/v2/postgres"
)

func (dao *privilegeLevelsDAO) Table() interfaces.IPostgresTable {
	return table.PrivilegeLevels
}

func (dao *privilegeLevelsDAO) InsertCols() postgres.ColumnList {
	return table.PrivilegeLevels.AllColumns.Except(
		table.PrivilegeLevels.ID,
		table.PrivilegeLevels.CreatedAt,
		table.PrivilegeLevels.UpdatedAt,
	)
}

func (dao *privilegeLevelsDAO) AllCols() postgres.ColumnList {
	return table.PrivilegeLevels.AllColumns
}

type privilegeLevelsDAO struct {
	db *sql.DB
}

func NewPrivilegeLevelsDAO(db *sql.DB) *privilegeLevelsDAO {
	return &privilegeLevelsDAO{
		db,
	}
}

func (dao *privilegeLevelsDAO) GetById(id int64) (*model.PrivilegeLevels, error) {
	var row model.PrivilegeLevels
	err := table.PrivilegeLevels.
		SELECT(table.PrivilegeLevels.AllColumns).
		FROM(table.PrivilegeLevels).
		WHERE(table.PrivilegeLevels.ID.EQ(postgres.Int(id))).
		LIMIT(1).
		Query(dao.db, &row)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (dao *privilegeLevelsDAO) Insert(row *model.PrivilegeLevels) error {
	var res model.PrivilegeLevels
	err := table.PrivilegeLevels.
		INSERT(table.PrivilegeLevels.Name).
		VALUES(postgres.String(row.Name)).
		RETURNING(table.PrivilegeLevels.ID).
		Query(dao.db, res)
	return err
}

func (dao *privilegeLevelsDAO) InsertMany(rows []*model.PrivilegeLevels) error {
	if len(rows) == 0 {
		return nil
	}
	return table.PrivilegeLevels.
		INSERT(table.PrivilegeLevels.Name).
		MODELS(rows).
		RETURNING(table.PrivilegeLevels.ID).
		Query(dao.db, &rows)
}

// Returns ID int64 if inserted.
// Parameter cols_update are the columns to be updated on conflict - if not provided, a few columns are updated
func (dao *privilegeLevelsDAO) Upsert(row *model.PrivilegeLevels, colsUpdate ...postgres.ColumnAssigment) error {
	if len(colsUpdate) == 0 {
		colsUpdate = []postgres.ColumnAssigment{
			table.PrivilegeLevels.Name.SET(postgres.String(row.Name)),
		}
	}

	row.UpdatedAt = tools.Ptr(time.Now())

	return table.PrivilegeLevels.
		INSERT(table.PrivilegeLevels.Name).
		VALUES(row.Name).
		ON_CONFLICT(table.PrivilegeLevels.Name).
		DO_UPDATE(postgres.SET(colsUpdate...)).
		RETURNING(table.PrivilegeLevels.ID).
		Query(dao.db, row)
}

func (dao *privilegeLevelsDAO) UpsertMany(rows []*model.PrivilegeLevels, colsUpdate ...postgres.ColumnAssigment) error {
	if len(colsUpdate) == 0 {
		colsUpdate = []postgres.ColumnAssigment{
			table.PrivilegeLevels.Name.SET(table.PrivilegeLevels.Name),
		}
	}

	now := time.Now()
	for _, r := range rows {
		r.UpdatedAt = &now
	}

	return table.PrivilegeLevels.
		INSERT(table.PrivilegeLevels.Name).
		MODELS(rows).
		ON_CONFLICT(table.PrivilegeLevels.Name).
		DO_UPDATE(postgres.SET(colsUpdate...)).
		RETURNING(table.PrivilegeLevels.ID).
		Query(dao.db, &rows)
}

func (dao *privilegeLevelsDAO) Update(row *model.PrivilegeLevels) error {
	row.UpdatedAt = tools.Ptr(time.Now())
	_, err := table.PrivilegeLevels.
		UPDATE(table.PrivilegeLevels.Name).
		MODEL(row).
		WHERE(table.PrivilegeLevels.ID.EQ(postgres.Int(row.ID))).
		SET(table.PrivilegeLevels.UpdatedAt.SET(postgres.TimestampT(time.Now()))).
		Exec(dao.db)
	return err
}

func (dao *privilegeLevelsDAO) Delete(id int64) error {
	_, err := table.PrivilegeLevels.
		DELETE().
		WHERE(table.PrivilegeLevels.ID.EQ(postgres.Int(id))).
		Exec(dao.db)
	if err != nil {
		return err
	}
	return nil
}

func (dao *privilegeLevelsDAO) GetAll() (*[]model.PrivilegeLevels, error) {
	var rows []model.PrivilegeLevels
	err := table.PrivilegeLevels.
		SELECT(table.PrivilegeLevels.AllColumns).
		ORDER_BY(table.PrivilegeLevels.ID.DESC()).
		Query(dao.db, &rows)
	if err != nil {
		return nil, err
	}
	return &rows, nil
}
