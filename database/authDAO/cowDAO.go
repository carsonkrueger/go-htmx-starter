package authDAO

import (
    "database/sql"
    "time"

    "github.com/carsonkrueger/main/gen/go_db/Cow/model"
    "github.com/carsonkrueger/main/gen/go_db/Cow/table"
    "github.com/carsonkrueger/main/tools"
    "github.com/go-jet/jet/v2/postgres"
)

type CowDAO struct {
    db *sql.DB
}

func NewCowDAO(db *sql.DB) *CowDAO {
    return &CowDAO{
        db,
    }
}

func (dao *CowDAO) GetById(id int64) (*model.Cow, error) {
    var row model.Cow
    err := table.Cow.SELECT(table.Cow.AllColumns).
        FROM(table.Cow).
        WHERE(table.Cow.ID.EQ(postgres.Int(id))).
        LIMIT(1).
        Query(dao.db, &row)
    if err != nil {
        return nil, err
    }
    return &row, nil
}

func (dao *CowDAO) Insert(row *model.Cow) error {
    var res model.Cow
    err := table.Cow.
        INSERT(table.Cow.AllColumns.Except(table.Cow.ID, table.Cow.CreatedAt, table.Cow.UpdatedAt)).
        VALUES(postgres.String(row.Name)).
        RETURNING(table.Cow.ID).
        Query(dao.db, res)
    return err
}

func (dao *CowDAO) InsertMany(rows []*model.Cow) error {
    if len(rows) == 0 {
        return nil
    }
    return table.Cow.
        INSERT(table.Cow.AllColumns.Except(table.Cow.ID, table.Cow.CreatedAt, table.Cow.UpdatedAt)).
        MODELS(rows).
        RETURNING(table.Cow.ID).
        Query(dao.db, &rows)
}

func (dao *CowDAO) Upsert(row *model.Cow, colsUpdate ...postgres.ColumnAssigment) error {
    if len(colsUpdate) == 0 {
        colsUpdate = []postgres.ColumnAssigment{
            table.Cow.Name.SET(postgres.String(row.Name)),
        }
    }

    row.UpdatedAt = tools.Ptr(time.Now())

    return table.Cow.
        INSERT(table.Cow.AllColumns.Except(table.Cow.ID, table.Cow.CreatedAt, table.Cow.UpdatedAt)).
        VALUES(row.Name).
        ON_CONFLICT(table.Cow.Name).
        DO_UPDATE(postgres.SET(colsUpdate...)).
        RETURNING(table.Cow.ID).
        Query(dao.db, row)
}

func (dao *CowDAO) UpsertMany(rows []*model.Cow, colsUpdate ...postgres.ColumnAssigment) error {
    if len(colsUpdate) == 0 {
        colsUpdate = []postgres.ColumnAssigment{
            table.Cow.Name.SET(table.Cow.Name),
        }
    }

    now := time.Now()
    for _, r := range rows {
        r.UpdatedAt = &now
    }

    return table.Cow.
        INSERT(table.Cow.AllColumns.Except(table.Cow.ID, table.Cow.CreatedAt, table.Cow.UpdatedAt)).
        MODELS(rows).
        ON_CONFLICT(table.Cow.Name).
        DO_UPDATE(postgres.SET(colsUpdate...)).
        RETURNING(table.Cow.ID).
        Query(dao.db, &rows)
}

func (dao *CowDAO) Update(row *model.Cow) error {
    row.UpdatedAt = tools.Ptr(time.Now())
    _, err := table.Cow.
        UPDATE(table.Cow.EXCLUDED.ID).
        MODEL(row).
        WHERE(table.Cow.ID.EQ(postgres.Int(row.ID))).
        SET(table.Cow.UpdatedAt.SET(postgres.TimestampT(time.Now()))).
        Exec(dao.db)
    return err
}

func (dao *CowDAO) Delete(id int64) error {
    _, err := table.Cow.DELETE().WHERE(table.Cow.ID.EQ(postgres.Int(id))).Exec(dao.db)
    if err != nil {
        return err
    }
    return nil
}

func (dao *CowDAO) GetAll() (*[]model.Cow, error) {
    var rows []model.Cow
    err := table.Cow.
        SELECT(table.Cow.AllColumns).
        ORDER_BY(table.Cow.ID.DESC()).
        Query(dao.db, &rows)
    if err != nil {
        return nil, err
    }
    return &rows, nil
}
