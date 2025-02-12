package repository

import (
	"database/sql"
	"errors"
	"exampleapp/domain/entity"
	"exampleapp/infrastructure/postgres"
	"fmt"
)

type ProductRepositoryImpl struct {
	conn *postgres.Conn
}

func NewProductRepositoryImpl(conn *postgres.Conn) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{conn}
}

func (r *ProductRepositoryImpl) Find(id uint64) *entity.Product {
	product := new(entity.Product)
	err := r.conn.QueryRow("select id, name, price from products where id = $1", id).Scan(&product.Id, &product.Name, &product.Price)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil
	case err != nil:
		panic(err)
	default:
		return product
	}
}

func (r *ProductRepositoryImpl) Paginate(page uint64, limit uint64) ([]*entity.Product, uint64) {
	var total uint64
	err := r.conn.QueryRow("select count(*) from products").Scan(&total)

	if err != nil {
		panic(err)
	}

	if total == 0 {
		return []*entity.Product{}, 0
	}

	rows, err := r.conn.Query(fmt.Sprintf("SELECT id, name, price from products limit %d offset %d", limit, (page-1)*limit))

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var products []*entity.Product

	for rows.Next() {
		product := new(entity.Product)

		if err := rows.Scan(&product.Id, &product.Name, &product.Price); err != nil {
			panic(err)
		}

		products = append(products, product)
	}

	return products, total
}
