package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/yosaZiege/jewelery-website/db"
	"github.com/yosaZiege/jewelery-website/models"
)




func AddProductToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx , cancel := context.WithTimeout(context.Background() , 100*time.Second)
		defer cancel()

		var cart models.Cart
		if err := c.BindJSON(&cart); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "invalid Input"})
			return
		}

		validationErr := validate.Struct(cart)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error" : validationErr.Error()})
			return
		}

		query := "INSERT INTO cart (user_id, product_id , quantity , created_at , updated_at ) VALUES($1,$2,$3, Now() , Now() );"

		_ , err := db.GetDB().ExecContext(ctx , query, cart.UserID , cart.ProductID , cart.Quantity)
		if err != nil {
			if pqErr , ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			 c.JSON(http.StatusConflict, gin.H{"error": "Product with this name already exists"})
			} else {
				fmt.Println("Error executing query:", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add product"})
			}
			return
		}	

		c.JSON(http.StatusOK , gin.H{"success" : "Product added to cart successfully"})
	} 

	}


func ViewCartProducts() gin.HandlerFunc {
	return func(c *gin.Context) {

		userId := c.Param("user_id")

		ctx, cancel := context.WithTimeout(context.Background() , 100*time.Second)
		defer cancel()

		var cartItems []models.Cart

		query := "SELECT product_id , quantity , updated_at from cart where user_id = $1"

		rows, err := db.GetDB().QueryContext(ctx ,query , userId)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"Error" : "Could not retrieve cart Items"})
			return
		}
		defer rows.Close()

		for rows.Next(){
			var cartItem models.Cart

			if err := rows.Scan(&cartItem.ProductID , &cartItem.Quantity , &cartItem.UpdatedAt); err != nil {
				fmt.Println(err)
				c.JSON(http.StatusInternalServerError , gin.H{"Error" : "Error scanning Items"})
				return
			}

			cartItems = append(cartItems, cartItem)
		}
		if err = rows.Err() ;err!= nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError  , gin.H{"Error" : "Error Retrieving Users"})
			return
		}

		c.JSON(http.StatusOK , cartItems)
	}
}
func UpdateQuantity() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        var cartItem models.Cart

        // Bind the JSON body to the cartItem variable
        if err := c.BindJSON(&cartItem); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request body"})
            return
        }

        // Update the quantity of the cart item in the database
        query := "UPDATE cart SET quantity = $1, updated_at = NOW() WHERE user_id = $2 AND product_id = $3 RETURNING id, user_id, product_id, quantity, created_at, updated_at;"

        var updatedCart models.Cart

        // Scan the result of the update into updatedCart
        err := db.GetDB().QueryRowContext(ctx, query, cartItem.Quantity, cartItem.UserID, cartItem.ProductID).Scan(&updatedCart.ID, &updatedCart.UserID, &updatedCart.ProductID, &updatedCart.Quantity, &updatedCart.AddedAt, &updatedCart.UpdatedAt)
        if err != nil {
            fmt.Println(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating the product"})
            return
        }

        // Return the updated cart item
        c.JSON(http.StatusOK, gin.H{"success": updatedCart})
    }
}

func RemoveProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()


		userID := c.Query("user_id")
        productID := c.Query("product_id")

		if userID == "" || productID == "" {
			c.JSON(http.StatusInternalServerError , gin.H{"error" : "User id or Product Id not founc" })
		}
		
		query := "DELETE FROM cart WHERE user_id =$1 and product_id=$2 RETURNING id"

		_ , err := db.GetDB().ExecContext(ctx , query , userID  , productID)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError , gin.H{"error" : "Failed to Delete Item from Cart"})
			return
		}

		c.JSON(http.StatusOK , gin.H{"success" :"Item Removed Successfully"})
		
	}
}
func GetCartItemsByUserId(userId int) ([]models.Cart, error){

	var cartItems []models.Cart

	query := "SELECT product_id, quantity FROM cart WHERE user_id = $1;"
    rows , err := db.GetDB().Query(query , userId)
	if err != nil {
       return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cartItem models.Cart

		if err := rows.Scan(&cartItem.ProductID , &cartItem.Quantity); err != nil {
			return nil ,err
		}
		cartItems = append(cartItems, cartItem)
	}
	return cartItems , nil
}