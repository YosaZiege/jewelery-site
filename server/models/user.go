package models

import (
	"time"
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID                 int       `json:"id" db:"id"`
	Username           string    `json:"username" db:"username" validate:"required,min=3,max=30"` // Ensure valid username length
	Email              string    `json:"email" db:"email" validate:"required,email"`
	Password           string    `json:"password" db:"password" validate:"required,min=8"` // Ensure strong password
	Role               string    `json:"role" db:"role" validate:"required"`
	IsEmailVerified    bool      `json:"is_email_verified" db:"is_email_verified"`
	PasswordResetToken string    `json:"password_reset_token" db:"password_reset_token"`
	Token              string    `json:"token"`
	RefreshToken       string    `json:"refresh_token"`
	ImageUrl           string    `json:"image_url"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}
var Validate *validator.Validate
func init() {
	Validate = validator.New() // Initialize the validator
}