package order

import "time"

type Order struct {
	ID         int       `json:"id" db:"id"`
	ProductID  int       `json:"productId" db:"product_id"`
	TotalPrice float64   `json:"totalPrice" db:"total_price"`
	Status     string    `json:"status" db:"status"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}
