package dao

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const USER_TABLE_NAME = "user_data"
const USER_CREATE_COLS = "email, salt, password_hash"
const USER_SELECT_COLS = "user_id, email"
const USER_PK_COL = "user_id"

type User struct {
	UserID       int    `json:"id" db:"user_id"`
	Email        string `json:"email"`
	Salt         []byte `json:"-"`
	PasswordHash []byte `json:"-" db:"password_hash"`
}

type UserDAO struct {
	BaseDAO[User]
}

func hashPassword(password string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write(salt)
	return hash.Sum(nil)
}

func NewUserDAO(db *sqlx.DB) *UserDAO {
	return &UserDAO{
		BaseDAO[User]{
			db:              db,
			selectStatement: "SELECT " + USER_SELECT_COLS + " FROM " + USER_TABLE_NAME + " OFFSET $1 LIMIT $2",
			getStatement:    "SELECT " + USER_SELECT_COLS + " FROM " + USER_TABLE_NAME + " WHERE user_id=$1",
		},
	}
}

func (dao *UserDAO) Insert(user *User) (int, error) {
	sqlStatement := "INSERT INTO " + USER_TABLE_NAME + " (" + USER_CREATE_COLS + ") VALUES ($1, $2, $3) RETURNING " + USER_PK_COL
	fmt.Println(sqlStatement)
	id := 0
	err := dao.db.QueryRow(sqlStatement, user.Email, user.Salt, user.PasswordHash).Scan(&id)
	return id, err
}

func (dao *UserDAO) GetByEmail(email string) (*User, error) {
	const sqlStatement = "SELECT " + USER_SELECT_COLS + " FROM " + USER_TABLE_NAME + " WHERE email=$1 LIMIT 1"

	user := User{}
	row := dao.db.QueryRow(sqlStatement, email)
	err := row.Scan(&user.UserID, &user.Email, &user.Salt, &user.PasswordHash)

	return &user, err
}

func (dao *UserDAO) Login(email string, password string) (int, error) {
	user, err := dao.GetByEmail(email)
	if err != nil {
		return 0, err
	}

	if !bytes.Equal(hashPassword(password, user.Salt), user.PasswordHash) {
		return 0, errors.New("invalid password")
	}

	return user.UserID, nil
}

func (dao *UserDAO) Register(email string, password string) (int, error) {
	salt := make([]byte, 16)

	_, err := rand.Read(salt)
	if err != nil {
		return 0, errors.New("could not generate salt")
	}

	passwordHash := hashPassword(password, salt)

	return dao.Insert(&User{Email: email, Salt: salt, PasswordHash: passwordHash})
}
