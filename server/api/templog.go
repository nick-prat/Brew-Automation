package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"raspberrysour/dao"
)

func (env *RequestEnvironment) GetTempLog(_ http.ResponseWriter, r *http.Request) (string, error) {
	tempLogDAO := dao.NewTempLogDAO(env.db)
	pk, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return "", BadRequestError(err)
	}

	tempLog, err := tempLogDAO.Get(pk)
	if err != nil {
		return "", InternalServerError(err)
	}

	resp, err := json.Marshal(tempLog)
	if err != nil {
		return "", InternalServerError(err)
	}

	return string(resp), nil
}

func (env *RequestEnvironment) GetTempLogs(_ http.ResponseWriter, _ *http.Request) (string, error) {
	tempLogDAO := dao.NewTempLogDAO(env.db)

	tempLogs, err := tempLogDAO.Select(10)
	if err != nil {
		return "", InternalServerError(err)
	}

	resp, err := json.Marshal(tempLogs)
	if err != nil {
		return "", InternalServerError(err)
	}

	return string(resp), nil
}

func (env *RequestEnvironment) PostTempLog(_ http.ResponseWriter, r *http.Request) (string, error) {
	tempLogDAO := dao.NewTempLogDAO(env.db)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", BadRequestError(err)
	}

	tempLog := dao.TempLog{}
	err = json.Unmarshal(body, &tempLog)
	if err != nil {
		return "", BadRequestError(err)
	}

	pk, err := tempLogDAO.Create(&tempLog)
	if err != nil {
		return "", InternalServerError(err)
	}

	resp, err := PKResponse(pk)
	if err != nil {
		return "", InternalServerError(err)
	}

	return resp, nil
}
