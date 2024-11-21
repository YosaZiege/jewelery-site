package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Order struct {
	ID            int       `json:"id" db:"id"`
	UserID        int       `json:"user_id" db:"user_id"`
	OrderDate     time.Time `json:"order_date" db:"order_date"`
	TotalAmount   float64   `json:"total_amount" db:"total_amount" validate:"required,min=0"` // Ensure total amount is not negative
	PaymentStatus string    `json:"payment_status" db:"payment_status" validate:"required"` // Ensure payment status is provided
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
func init() {
	Validate = validator.New()
}