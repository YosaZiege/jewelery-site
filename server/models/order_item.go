package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type OrderItem struct {
    ID        int       `json:"id" db:"id"`
    OrderID   int       `json:"order_id" db:"order_id"`
    ProductID int       `json:"product_id" db:"product_id"`
    Quantity  int       `json:"quantity" db:"quantity"`
    Price     float64   `json:"price" db:"price"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func init() {
	Validate = validator.New()
}