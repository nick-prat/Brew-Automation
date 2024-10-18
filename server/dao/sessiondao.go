package dao

import (
	"database/sql"
	"encoding/base64"
)

const SESSION_TABLE_NAME = "session"
const SESSION_CREATE_COLS = "user_id"
const SESSION_SELECT_COLS = "session_id, user_id"
const SESSION_PK_COL = "session_id"

type SessionDAO struct {
	db *sql.DB
}

type Session struct {
	SessionId int
	UserId    int
	Token     string
}

func NewSessionDAO(db *sql.DB) *SessionDAO {
	return &SessionDAO{db: db}
}

func (dao *SessionDAO) Get(token string) (*Session, error) {
	const sqlStatement = "SELECT " + SESSION_SELECT_COLS + " FROM " + SESSION_TABLE_NAME + " WHERE token=$1 LIMIT 1"

	bytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	var userId int
	var sessionId int
	row := dao.db.QueryRow(sqlStatement, bytes)
	err = row.Scan(&sessionId, &userId)

	switch err {
	case nil:
		return &Session{SessionId: sessionId, UserId: userId, Token: token}, nil
	default:
		return nil, err
	}
}

func (dao *SessionDAO) Create(userId int) (string, error) {
	sqlStatement := "INSERT INTO " + SESSION_TABLE_NAME + " (" + SESSION_CREATE_COLS + ") VALUES ($1) RETURNING " + "token"
	bytes := make([]byte, 0)
	err := dao.db.QueryRow(sqlStatement, userId).Scan(&bytes)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}
