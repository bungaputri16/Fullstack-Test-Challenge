package order

import (
	"database/sql"
	"strconv"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(o *Order) error {
	query := `INSERT INTO orders (product_id, total_price, status, created_at) 
	          VALUES ($1, $2, $3, NOW()) RETURNING id`
	return r.db.QueryRow(query, o.ProductID, o.TotalPrice, o.Status).
		Scan(&o.ID)
}

func (r *Repository) FindByProductID(productId int) ([]Order, error) {
	rows, err := r.db.Query(
		"SELECT id, product_id, total_price, status, created_at FROM orders WHERE product_id=$1", 
		productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

var orders []Order
	for rows.Next() {
		var o Order
		var totalPriceStr string // ambil total_price sebagai string

		if err := rows.Scan(&o.ID, &o.ProductID, &totalPriceStr, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}

		// convert string → float64
		price, err := strconv.ParseFloat(totalPriceStr, 64)
		if err != nil {
			price = 0
		}
		o.TotalPrice = price

		orders = append(orders, o)
	}
	return orders, nil
}