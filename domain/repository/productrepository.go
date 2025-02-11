package repository

import "exampleapp/domain/entity"

type ProductRepository interface {
	//All() entity.ProductCollection
	Find(id uint64) *entity.Product
}
