package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
)

const USER_TABLE_NAME = "user_data"
const USER_CREATE_COLS = "email, salt, password_hash"
const USER_SELECT_COLS = "email, salt, password_hash"
const USER_PK_COL = "user_id"

type UserDAO struct {
	db *sql.DB
}

type User struct {
	UserID       int
	Email        string
	Salt         []byte
	PasswordHash []byte
}

type UserLogin struct {
	Email    string
	Password string
}

type UserRegister = UserLogin

func hashPassword(password string, salt []byte) []byte {
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write(salt)
	return hash.Sum(nil)
}

func NewUserDAO(db *sql.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (dao *UserDAO) Insert(user *User) (int, error) {
	sqlStatement := "INSERT INTO " + USER_TABLE_NAME + " (" + USER_CREATE_COLS + ") VALUES ($1, $2, $3) RETURNING " + USER_PK_COL
	fmt.Println(sqlStatement)
	id := 0
	err := dao.db.QueryRow(sqlStatement, user.Email, user.Salt, user.PasswordHash).Scan(&id)
	return id, err
}

func (dao *UserDAO) GetByPK(pk int) (*User, error) {
	const sqlStatement = "SELECT " + USER_SELECT_COLS + " FROM " + USER_TABLE_NAME + " WHERE user_id=$1 LIMIT 1"

	user := User{}
	row := dao.db.QueryRow(sqlStatement, pk)
	err := row.Scan(&user.Email, &user.Salt, &user.PasswordHash)

	return &user, err
}

func (dao *UserDAO) GetByEmail(email string) (*User, error) {
	const sqlStatement = "SELECT " + USER_SELECT_COLS + " FROM " + USER_TABLE_NAME + " WHERE email=$1 LIMIT 1"

	user := User{}
	row := dao.db.QueryRow(sqlStatement, email)
	err := row.Scan(&user.Email, &user.Salt, &user.PasswordHash)

	return &user, err
}

func (dao *UserDAO) Login(email string, password string) (*User, error) {
	user, err := dao.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(hashPassword(password, user.Salt), user.PasswordHash) {
		return nil, errors.New("invalid password")
	}

	return user, nil
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
