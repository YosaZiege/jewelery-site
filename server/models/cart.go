package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)





type Cart struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	ProductID int       `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity" validate:"required,min=1"` // Ensure quantity is at least 1
	AddedAt   time.Time `json:"added_at" db:"added_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func init() {
	Validate = validator.New()
}