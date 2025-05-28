package usecase

import (
	"context"
	"exampleapp/internal/domain/entity"
	"exampleapp/internal/domain/repository"
)

type PaginateObjectsUseCase struct {
	repository repository.ObjectRepository
}

func NewPaginateObjectsUseCase(
	repository repository.ObjectRepository,
) *PaginateObjectsUseCase {
	return &PaginateObjectsUseCase{repository}
}

func (u *PaginateObjectsUseCase) Handle(ctx context.Context, query repository.ObjectsQuery) *PaginateObjModel {
	objects, total := u.repository.Paginate(ctx, query)

	return &PaginateObjModel{
		Page:    query.Page,
		Limit:   query.Limit,
		Total:   total,
		Objects: objects,
		Metro:   u.repository.Metro(ctx),
		Rooms:   u.repository.Rooms(ctx),
	}
}

type PaginateObjModel struct {
	Page    uint64
	Limit   uint64
	Total   uint64
	Objects []*entity.Object
	Metro   []string
	Rooms   []uint8
}
