package dao

import (
	"time"

	"github.com/jmoiron/sqlx"
)

const FERMENT_RUN_TABLE_NAME string = "ferment_run"
const FERMENT_RUN_PK_COL string = "ferment_run_id"
const FERMENT_RUN_SELECT_COLUMN string = "ferment_run_id, name, start_date"

type FermentRun struct {
	Id        int       `json:"id" db:"ferment_run_id"`
	Name      float32   `json:"name" db:"name"`
	StartDate time.Time `json:"start_date" db:"start_date"`
}

type FermentRunDAO struct {
	createStatement string
	BaseDAO[FermentRun]
}

func NewFermentRunDAO(db *sqlx.DB) *FermentRunDAO {
	return &FermentRunDAO{
		createStatement: "INSERT INTO " + FERMENT_RUN_TABLE_NAME + " (name) VALUES ($1) RETURNING " + FERMENT_RUN_PK_COL,
		BaseDAO: BaseDAO[FermentRun]{
			db:              db,
			selectStatement: "SELECT " + FERMENT_RUN_SELECT_COLUMN + " FROM " + FERMENT_RUN_TABLE_NAME + " OFFSET $1 LIMIT $2",
			getStatement:    "SELECT " + FERMENT_RUN_SELECT_COLUMN + " FROM " + FERMENT_RUN_TABLE_NAME + " WHERE " + FERMENT_RUN_PK_COL + "=$1;",
		},
	}
}

func (dao *FermentRunDAO) Create(fermentRun *FermentRun) (int, error) {
	id := 0
	err := dao.db.QueryRow(dao.createStatement, fermentRun.Name).Scan(&id)
	return id, err
}
