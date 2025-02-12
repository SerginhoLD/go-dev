package repository

import "exampleapp/domain/entity"

type ProductRepository interface {
	Find(id uint64) *entity.Product
	Paginate(page uint64, limit uint64) ([]*entity.Product, uint64)
}
