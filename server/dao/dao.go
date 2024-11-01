package dao

import "github.com/jmoiron/sqlx"

type DAO[T any] interface {
	Select(limit int) ([]*T, error)
	Create(*T) (int, error)
	Get(pk int) (*T, error)
}

type BaseDAO[T any] struct {
	db              *sqlx.DB
	selectStatement string
	getStatement    string
}

func (dao *BaseDAO[T]) Select(skip int, limit int) (*[]T, error) {
	rows := []T{}
	err := dao.db.Select(&rows, dao.selectStatement, skip, limit)
	if err != nil {
		return nil, err
	}
	return &rows, nil
}

func (dao *BaseDAO[T]) Get(pk int) (*T, error) {
	row := new(T)
	err := dao.db.Get(row, dao.getStatement, pk)
	if err != nil {
		return nil, err
	}
	return row, nil
}
