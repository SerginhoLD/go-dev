package repository

import (
	"database/sql"
	"errors"
	"exampleapp/domain/entity"
)

type ProductRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepositoryImpl(db *sql.DB) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{db}
}

func (r *ProductRepositoryImpl) Find(id uint64) *entity.Product {
	product := new(entity.Product)
	err := r.db.QueryRow("select id, name, price from products where id = $1", id).Scan(&product.Id, &product.Name, &product.Price)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil
	case err != nil:
		panic(err)
	default:
		return product
	}
}

func (r *ProductRepositoryImpl) All() *entity.ProductCollection {
	var count uint64
	err := r.db.QueryRow("select count(*) from products").Scan(&count)

	if err != nil {
		panic(err)
	}

	rows, err := r.db.Query("SELECT id, name, price from products")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	collection := &entity.ProductCollection{count, make([]*entity.Product, 0)}

	for rows.Next() {
		product := new(entity.Product)

		if err := rows.Scan(&product.Id, &product.Name, &product.Price); err != nil {
			panic(err)
		}

		collection.Products = append(collection.Products, product)
	}

	return collection
}
