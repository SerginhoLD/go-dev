package repository

import (
	"context"
	"exampleapp/domain/entity"
)

type ProductRepository interface {
	Find(ctx context.Context, id uint64) *entity.Product
	Paginate(ctx context.Context, page uint64, limit uint64) ([]*entity.Product, uint64)
}
