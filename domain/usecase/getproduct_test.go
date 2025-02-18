package usecase

import (
	"context"
	"exampleapp/domain/entity"
	"fmt"
	"testing"
)

func TestGetProductUseCase(t *testing.T) {
	var tests = []struct {
		id  uint64
		has bool
	}{
		{id: 0, has: false},
		{id: 1, has: true},
		{id: 2, has: true},
		{id: 3, has: false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("id: %d", tt.id), func(t *testing.T) {
			useCase := NewGetProductUseCase(&MockRepository{})
			product := useCase.Handle(context.Background(), GetProductQuery{Id: tt.id})

			if tt.has && product.Id != tt.id {
				t.Errorf("got %d, want %d", product.Id, tt.id)
			}

			if !tt.has && product != nil {
				t.Errorf("got %T, want nil", product.Id)
			}
		})
	}
}

type MockRepository struct {
}

func (m *MockRepository) Find(ctx context.Context, id uint64) *entity.Product {
	if id == 0 || id > 2 {
		return nil
	}

	return &entity.Product{
		Id: id,
	}
}

func (m *MockRepository) FindByName(ctx context.Context, name string) *entity.Product {
	panic("not tested")
}

func (m *MockRepository) Paginate(ctx context.Context, page uint64, limit uint64) ([]*entity.Product, uint64) {
	panic("not tested")
}

func (m *MockRepository) Create(ctx context.Context, product *entity.Product) {
	panic("not tested")
}
