package repository

import (
	"context"
	"exampleapp/internal/domain/entity"
)

type ObjectRepository interface {
	Paginate(ctx context.Context, query ObjectsQuery) (objects []*entity.Object, total uint64)
	Save(ctx context.Context, objects []*entity.Object)
	Metro(ctx context.Context) []string
	Rooms(ctx context.Context) []uint8
}

type ObjectsQuery struct {
	Limit     uint64
	Page      uint64 `schema:"page,default=1"`
	Metro     string `schema:"metro"`
	Loc       uint8  `schema:"loc"`
	Search    string `schema:"q"`
	PriceFrom uint64 `schema:"min"`
	PriceTo   uint64 `schema:"max"`
	SizeFrom  uint64 `schema:"sizeFrom"`
	SizeTo    uint64 `schema:"sizeTo"`
	Rooms     uint8  `schema:"rooms"`
	Checked   int8   `schema:"checked"` // -1 - нет, 1 - да, 0 - все
}
