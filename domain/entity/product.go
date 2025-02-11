package entity

type Product struct {
	Id    uint64
	Name  string
	Price float64
}

type ProductCollection struct {
	Total    uint64
	Products []Product
}
