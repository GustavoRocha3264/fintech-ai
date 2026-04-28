package portfolio

import "errors"

var ErrNotFound = errors.New("portfolio not found")

type Repository interface {
	Save(p Portfolio) error
	FindByID(id string) (*Portfolio, error)
}
