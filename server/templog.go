package main

import (
	"database/sql"
	"raspberrysour/dao"
)

type TempLog struct {
	Id         int     `json:"id"`
	Temp       float32 `json:"temp"`
	FermentRun int     `json:"ferment_run_id"`
}

type TempLogDAO struct {
	db *sql.DB
}

func NewTempLogDAO(db *sql.DB) dao.DAO[TempLog] {
	return &TempLogDAO{db: db}
}

func (dao *TempLogDAO) Create(tempLog *TempLog) (int, error) {
	sqlStatement := "INSERT INTO temp_log (temp, ferment_run_id) VALUES ($1, $2) RETURNING temp_log_id"
	id := 0
	err := dao.db.QueryRow(sqlStatement, tempLog.Temp, tempLog.FermentRun).Scan(&id)
	return id, err
}

func (dao *TempLogDAO) Select(limit int) ([]*TempLog, error) {
	const sqlStatement = "SELECT temp_log_id, ferment_run_id, temp FROM temp_log LIMIT $1"

	rows, err := dao.db.Query(sqlStatement, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	logs := []*TempLog{}
	for rows.Next() {
		log := TempLog{}
		err := rows.Scan(&log.Id, &log.FermentRun, &log.Temp)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (dao *TempLogDAO) Get(pk int) (*TempLog, error) {
	sqlStatement := "SELECT temp_log_id, ferment_run_id, temp FROM temp_log WHERE temp_log_id=$1;"

	var tempLogId int
	var fermentRunId int
	var temp float32

	row := dao.db.QueryRow(sqlStatement, pk)
	err := row.Scan(&tempLogId, &fermentRunId, &temp)

	switch err {
	case nil:
		return &TempLog{Id: tempLogId, Temp: temp, FermentRun: fermentRunId}, nil
	default:
		return nil, err
	}

}
