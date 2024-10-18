package dao

type DAO[T any] interface {
	Select(limit int) ([]*T, error)
	Create(*T) (int, error)
	Get(pk int) (*T, error)
}
