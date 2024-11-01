package dao

import (
	"time"

	"github.com/jmoiron/sqlx"
)

const TEMP_LOG_TABLE_NAME string = "temp_log"
const TEMP_LOG_SELECT_COLUMSN string = "temp_log_id, ferment_run_id, temp, time_stamp"

type TempLog struct {
	Id         int       `json:"id" db:"temp_log_id"`
	Temp       float32   `json:"temp" db:"temp"`
	FermentRun int       `json:"ferment_run_id" db:"ferment_run_id"`
	TimeStamp  time.Time `json:"time_stamp" db:"time_stamp"`
}

type TempLogDAO struct {
	BaseDAO[TempLog]
}

func NewTempLogDAO(db *sqlx.DB) *TempLogDAO {
	return &TempLogDAO{
		BaseDAO[TempLog]{
			db:              db,
			selectStatement: "SELECT " + TEMP_LOG_SELECT_COLUMSN + " FROM " + TEMP_LOG_TABLE_NAME + " OFFSET $1 LIMIT $2",
			getStatement:    "SELECT " + TEMP_LOG_SELECT_COLUMSN + " FROM " + TEMP_LOG_TABLE_NAME + " WHERE temp_log_id=$1;",
		},
	}
}

func (dao *TempLogDAO) Create(tempLog *TempLog) (int, error) {
	sqlStatement := "INSERT INTO temp_log (temp, ferment_run_id) VALUES ($1, $2) RETURNING temp_log_id"
	id := 0
	err := dao.db.QueryRow(sqlStatement, tempLog.Temp, tempLog.FermentRun).Scan(&id)
	return id, err
}
