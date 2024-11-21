package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type ShippingDetails struct {
	ID                 int       `json:"id" db:"id"`
	OrderID            int       `json:"order_id" db:"order_id" validate:"required"` // Ensure order ID is provided
	ShippingName       string    `json:"shipping_name" db:"shipping_name" validate:"required"`
	ShippingStreet     string    `json:"shipping_street" db:"shipping_street" validate:"required"`
	ShippingCity       string    `json:"shipping_city" db:"shipping_city" validate:"required"`
	ShippingState      string    `json:"shipping_state" db:"shipping_state" validate:"required"`
	ShippingPostalCode string    `json:"shipping_postal_code" db:"shipping_postal_code" validate:"required"`
	ShippingCountry    string    `json:"shipping_country" db:"shipping_country" validate:"required"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}
func init() {
	Validate = validator.New()
}