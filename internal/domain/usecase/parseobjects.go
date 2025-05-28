package usecase

import (
	"context"
	"exampleapp/internal/domain/messages"
	"exampleapp/internal/domain/repository"
)

type ParseObjectsUseCase struct {
	objRepository   repository.ObjectRepository
	veObjRepository repository.VeObjectRepository
}

func NewParseObjectsUseCase(
	objRepository repository.ObjectRepository,
	veObjRepository repository.VeObjectRepository,
) *ParseObjectsUseCase {
	return &ParseObjectsUseCase{
		objRepository:   objRepository,
		veObjRepository: veObjRepository,
	}
}

func (u *ParseObjectsUseCase) Handle(ctx context.Context, command *messages.ParsePageMessage) error {
	veObjects := u.veObjRepository.Paginate(ctx, command.Page)

	if len(veObjects) == 0 {
		return nil
	}

	u.objRepository.Save(ctx, veObjects)

	return nil
}
