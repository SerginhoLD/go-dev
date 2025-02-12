package repository

import "exampleapp/domain/entity"

type ProductRepository interface {
	Find(id uint64) *entity.Product
	All() ([]*entity.Product, uint64)
}
