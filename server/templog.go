package main

import "database/sql"

type TempLog struct {
	Id         int     `json:"id"`
	Temp       float32 `json:"temp"`
	FermentRun int     `json:"ferment_run_id"`
}

type TempLogDAO struct {
	db *sql.DB
}

func NewTempLogDAO(db *sql.DB) DAO[TempLog] {
	return &TempLogDAO{db: db}
}

func (dao *TempLogDAO) Create(tempLog *TempLog) (int, error) {
	sqlStatement := "INSERT INTO temp_log (temp, ferment_run_id) VALUES ($1, $2) RETURNING temp_log_id"
	id := 0
	err := dao.db.QueryRow(sqlStatement, tempLog.Temp, tempLog.FermentRun).Scan(&id)
	return id, err
}

func (dao *TempLogDAO) Select(limit int) ([]*TempLog, error) {
	logs := []*TempLog{}
	logs = append(logs, &TempLog{})
	return logs, nil
}

func (dao *TempLogDAO) Get(pk int) (*TempLog, error) {
	sqlStatement := "SELECT temp_log_id, ferment_run_id, temp FROM temp_log WHERE id=$1;"

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
