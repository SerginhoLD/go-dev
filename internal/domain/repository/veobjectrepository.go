package repository

import (
	"context"
	"exampleapp/internal/domain/entity"
)

type VeObjectRepository interface {
	Total(ctx context.Context) (total uint64, size uint8)
	Paginate(ctx context.Context, page uint64) []*entity.Object
}
