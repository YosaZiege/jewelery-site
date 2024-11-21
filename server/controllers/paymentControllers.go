package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)


func init() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}

func CreatePaymentIntent(c *gin.Context) {
	var req struct {
		Amount int64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	params := &stripe.PaymentIntentParams{
		Amount: stripe.Int64(req.Amount),
		Currency: stripe.String("usd"),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Unable to create Payment Intent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"client_secret": pi.ClientSecret,
	})
}
