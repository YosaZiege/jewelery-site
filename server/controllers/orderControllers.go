package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yosaZiege/jewelery-website/db"
	"github.com/yosaZiege/jewelery-website/models"
)

func CreateOrderApi() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var input struct {
			Order          models.Order          `json:"order"`
			ShippingDetails models.ShippingDetails `json:"shipping_details"`
		}

		// Bind JSON body to both order and shipping details
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
			return
		}

		order := input.Order
		shippingDetails := input.ShippingDetails

		// Validate the order data
		validationErr := validate.Struct(order)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Insert the order into the orders table
		orderQuery := "INSERT INTO orders (user_id, total_amount, payment_status) VALUES ($1, $2, $3) RETURNING id;"
		err := db.GetDB().QueryRowContext(ctx, orderQuery, order.UserID, order.TotalAmount, order.PaymentStatus).Scan(&order.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create order"})
			return
		}

		// Set the OrderID in shipping details
		shippingDetails.OrderID = order.ID

		// Insert the shipping details into the shippingdetails table
		shippingQuery := "INSERT INTO shippingdetails (order_id, shipping_name, shipping_street, shipping_city, shipping_state, shipping_postal_code, shipping_country) VALUES ($1, $2, $3, $4, $5, $6, $7);"
		_, err = db.GetDB().ExecContext(ctx, shippingQuery, shippingDetails.OrderID, shippingDetails.ShippingName, shippingDetails.ShippingStreet, shippingDetails.ShippingCity, shippingDetails.ShippingState, shippingDetails.ShippingPostalCode, shippingDetails.ShippingCountry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create shipping details"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"success": order})
	}
}

func CreateOrder(userId int, totalAmount float64, shippingDetails models.ShippingDetails, paymentStatus string) (int, error) {
	var orderID int

	// Insert the order into the orders table
	query := "INSERT INTO orders (user_id, total_amount, payment_status) VALUES ($1, $2, $3) RETURNING id;"
	err := db.GetDB().QueryRow(query, userId, totalAmount, paymentStatus).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("could not create order: %w", err)
	}

	// Set the OrderID in shipping details
	shippingDetails.OrderID = orderID

	// Insert the shipping details into the shippingdetails table
	shippingQuery := "INSERT INTO shippingdetails (order_id, shipping_name, shipping_street, shipping_city, shipping_state, shipping_postal_code, shipping_country) VALUES ($1, $2, $3, $4, $5, $6, $7);"
	_, err = db.GetDB().Exec(shippingQuery, shippingDetails.OrderID, shippingDetails.ShippingName, shippingDetails.ShippingStreet, shippingDetails.ShippingCity, shippingDetails.ShippingState, shippingDetails.ShippingPostalCode, shippingDetails.ShippingCountry)
	if err != nil {
		return 0, fmt.Errorf("could not create shipping details: %w", err)
	}

	return orderID, nil
}

func GetOrderById() gin.HandlerFunc {
    return func(c *gin.Context) {
        orderID := c.Param("order_id")
        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        var order models.Order
        query := "SELECT * FROM orders WHERE id = $1;"
        err := db.GetDB().QueryRowContext(ctx, query, orderID).Scan(&order.ID, &order.UserID, &order.OrderDate, &order.TotalAmount, &order.PaymentStatus, &order.CreatedAt, &order.UpdatedAt)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"order": order})
    }
}

func UpdateOrder() gin.HandlerFunc {
    return func(c *gin.Context) {
        orderID := c.Param("order_id")
        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        var order models.Order
        if err := c.BindJSON(&order); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request body"})
            return
        }

        query := "UPDATE orders SET total_amount = $1, payment_status = $2, updated_at = NOW() WHERE id = $3 RETURNING *;"
        var updatedOrder models.Order
        err := db.GetDB().QueryRowContext(ctx, query, order.TotalAmount, order.PaymentStatus, orderID).Scan(&updatedOrder.ID, &updatedOrder.UserID, &updatedOrder.OrderDate, &updatedOrder.TotalAmount, &updatedOrder.PaymentStatus, &updatedOrder.CreatedAt, &updatedOrder.UpdatedAt)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating the order"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"success": updatedOrder})
    }
}

func DeleteOrder() gin.HandlerFunc {
    return func(c *gin.Context) {
        orderID := c.Param("order_id")
        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        query := "DELETE FROM orders WHERE id = $1;"
        _, err := db.GetDB().ExecContext(ctx, query, orderID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting the order"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"success": "Order deleted successfully"})
    }
}


func AddOrderItem(orderId int,productId int,quantity int,price float64) error{
   
	query := "INSERT INTO order_items (order_id, product_id, quantity, price, created_at, updated_at) VALUES($1, $2, $3, $4, NOW(), NOW());"
	_ , err := db.GetDB().Exec(query , orderId , productId , quantity , price)
	if err != nil {
		return fmt.Errorf("could not add order item: %w",err)
	}
	return nil
}
func ClearCart(userId int) error {
	query := "DELETE FROM cart WHERE user_id = $1;"
	_ , err := db.GetDB().Exec(query , userId)
	if err != nil {
		return fmt.Errorf("could not clear Cart : %w", err)
	}
	return nil
}



func CreateOrderFromCart(userId int, shippingDetails models.ShippingDetails) error {
    cartItems, err := GetCartItemsByUserId(userId)
    if err != nil {
        return err
    }

    // Calculate total amount from cart items
    totalAmount := CalculateTotalAmount(cartItems)

    // Create an order
    orderId, err := CreateOrder(userId, totalAmount, shippingDetails, "pending")
    if err != nil {
        return err
    }

    // Add cart items to order_items
    for _, item := range cartItems {
        // Fetch price from the product table for each item
        product, err := GetProductById(item.ProductID)
        if err != nil {
            return fmt.Errorf("failed to fetch product for order item: %w", err)
        }
        // Pass product price to AddOrderItem
        if err := AddOrderItem(orderId, item.ProductID, item.Quantity, product.Price); err != nil {
            return fmt.Errorf("failed to add order item: %w", err)
        }
    }

    // Clear cart after order creation
    if err := ClearCart(userId); err != nil {
        return fmt.Errorf("failed to clear cart: %w", err)
    }

    return nil
}

func CalculateTotalAmount(cartItems []models.Cart) float64 {
    total := 0.0
    for _, item := range cartItems {
        var price float64
        query := "SELECT price FROM products WHERE id = $1;"
        err := db.GetDB().QueryRow(query, item.ProductID).Scan(&price)
        if err == nil {
            total += price * float64(item.Quantity)
        }
    }
    return total
}
