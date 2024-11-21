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

func GetAllProducts() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        var products []models.Product
        query := "SELECT id,name,description,price,discount,image_url FROM products;"

        rows, err := db.GetDB().QueryContext(ctx, query)
        if err != nil {
            fmt.Println(err)
            c.JSON(http.StatusInternalServerError , gin.H{"Error" : "Could not retrieve Users "})
            return
        }
        defer rows.Close()

        for rows.Next() {
            var product models.Product
            
            if err := rows.Scan(&product.ID ,&product.Name, &product.Description , &product.Price  , &product.Discount,&product.ImageUrl); err != nil {
                fmt.Println(err)
                c.JSON(http.StatusInternalServerError , gin.H{"Error" : "Scanning Error of Products"})
                return
            }
        products = append(products, product)
        }
        if err = rows.Err() ; err != nil {
            fmt.Println(err)
            c.JSON(http.StatusInternalServerError , gin.H{"Error" : "Error Retrieving Users"})
            return
        }

        c.JSON(http.StatusOK , products)
    }
}
func GetProductByIdApi() gin.HandlerFunc{
    return func(c *gin.Context) {
        productId := c.Param("product_id")

        ctx , cancel := context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        var product models.Product
        query := "SELECT id, name , description , price , discount , image_url FROM products where id = $1;"

        err := db.GetDB().QueryRowContext(ctx, query , productId).Scan(&product.ID,&product.Name, &product.Description , &product.Price  , &product.Discount, &product.ImageUrl)
        fmt.Println(err)
        if err != nil {
            c.JSON(http.StatusNotFound , gin.H{"Error" : "Product not found"})
            return
        }

        c.JSON(http.StatusOK , product)
    }
}
func GetProductById(productId int) (models.Product, error) {
    var product models.Product
    query := "SELECT id, name, description, price, discount, image_url,sold FROM products WHERE id = $1;"

    err := db.GetDB().QueryRow(query, productId).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Discount, &product.ImageUrl ,&product.Sold)
    
    if err != nil {
        fmt.Println(err)
        return product, fmt.Errorf("could not get product: %w", err)
    }
    return product, nil
}

func UpdateProduct() gin.HandlerFunc {
    return func(c *gin.Context) {
        productId := c.Param("product_id")
        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()
        
        var product models.Product

        if err := c.BindJSON(&product); err != nil {
            fmt.Println(err)
            c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request"})
            return
        }

        query := "UPDATE products SET name=$1, description=$2, price=$3, discount=$4, image_url=$5, sold=$6, updated_at=NOW() WHERE id=$7 RETURNING id, name, description, price, discount, image_url, created_at, updated_at, sold;"

        var updatedProduct models.Product
        err := db.GetDB().QueryRowContext(
            ctx, query,
            product.Name,
            product.Description,
            product.Price,
            product.Discount,
            product.ImageUrl,
            product.Sold,
            productId,
        ).Scan(
            &updatedProduct.ID,
            &updatedProduct.Name,
            &updatedProduct.Description,
            &updatedProduct.Price,
            &updatedProduct.Discount,
            &updatedProduct.ImageUrl,
            &updatedProduct.CreatedAt,
            &updatedProduct.UpdatedAt,
            &updatedProduct.Sold,
        )

        if err != nil {
            fmt.Println(err)
            c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error updating the product"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"success": updatedProduct})
    }
}

func AddProduct() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        var product models.Product

        if err := c.BindJSON(&product); err != nil {
            fmt.Println("Error binding JSON:", err)
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        validationErr := validate.Struct(product)
        if validationErr != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
            return
        }

        query := "INSERT INTO products (name, description, price, discount, image_url) VALUES ($1, $2, $3, $4, $5)"
        _, err := db.GetDB().ExecContext(ctx, query, product.Name, product.Description, product.Price, product.Discount, product.ImageUrl)
        if err != nil {
            if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
                c.JSON(http.StatusConflict, gin.H{"error": "Product with this name already exists"})
            } else {
                fmt.Println("Error executing query:", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add product"})
            }
            return
        }

        c.JSON(http.StatusCreated, gin.H{"message": "Product added successfully"})
    }
}

func DeleteProduct() gin.HandlerFunc {
    return func(c *gin.Context) {
        productId := c.Param("product_id")

        ctx , cancel := context.WithTimeout(context.Background() , 100*time.Second)
        defer cancel()

        query := "DELETE FROM products WHERE id=$1"
        
        _ , err := db.GetDB().ExecContext(ctx , query ,productId)
        if err != nil {
            c.JSON(http.StatusInternalServerError , gin.H{"Error" : "Failed to Delete Product"})
            return
        }

        c.JSON(http.StatusOK , gin.H{"success" : "Product Deleted Successfully "})
    }
}

func BestSellingProducts() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx , cancel := context.WithTimeout(context.Background() , 100*time.Second)
        defer cancel()

        query := "SELECT id, name , description,price  , discount , created_at , updated_at , image_url, sold FROM products ORDER BY sold DESC LIMIT 4;"
        var products []models.Product
        rows , err := db.GetDB().QueryContext(ctx, query)
        if err != nil {
            fmt.Println(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error":"Could not retrieve best selling products"})
            return
        }
        defer rows.Close()

        for rows.Next(){
            var product models.Product

            if err := rows.Scan(&product.ID, &product.Name , &product.Description , &product.Price , &product.Discount , &product.CreatedAt , &product.UpdatedAt , &product.ImageUrl ,&product.Sold); err != nil {
                fmt.Println((err))
                c.JSON(http.StatusInternalServerError, gin.H{"Error":"Error Scanning products"})
                return
            }
            products = append(products , product)
        }
        if err = rows.Err() ; err != nil{
            fmt.Println(err)
            c.JSON(http.StatusInternalServerError , gin.H{"Error" : "Error retreiving products"})
            return
        }

        c.JSON(http.StatusOK , products)
    }
}

func FetchProductPageDetails() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        productId := c.Param("product_id")

        // Structs for the data
        var productPage models.ProductPageData
        var colors []models.Color
        var reviews []models.Review
        var sizes []string
        var images []string
        var product models.Product

        // Query: Product details
        queryProduct := "SELECT name, description, price FROM products WHERE id = $1"
        errProduct := db.GetDB().QueryRowContext(ctx, queryProduct, productId).Scan(&product.Name, &product.Description, &product.Price)
        if errProduct != nil {
            fmt.Println(errProduct)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve product details"})
            return
        }

        // Query: Colors
        queryColors := "SELECT hex, color_name, color_image_url FROM product_colors WHERE product_id = $1"
        colorRows, errColor := db.GetDB().QueryContext(ctx, queryColors, productId)
        if errColor != nil {
            fmt.Println(errColor)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve colors"})
            return
        }
        defer colorRows.Close()

        for colorRows.Next() {
            var color models.Color
            if err := colorRows.Scan(&color.Hex , &color.Name, &color.ImageUrl); err != nil {
                fmt.Println(err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning colors"})
                return
            }
            colors = append(colors, color)
        }

        // Query: Reviews
        queryReviews := "SELECT r.rating, r.comment, u.username FROM reviews AS r JOIN users AS u ON u.id = r.user_id WHERE r.product_id = $1"
        reviewRows, errReviews := db.GetDB().QueryContext(ctx, queryReviews, productId)
        if errReviews != nil {
            fmt.Println(errReviews)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve reviews"})
            return
        }
        defer reviewRows.Close()

        for reviewRows.Next() {
            var review models.Review
            if err := reviewRows.Scan(&review.Rating, &review.Comment, &review.Username); err != nil {
                fmt.Println(err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning reviews"})
                return
            }
            reviews = append(reviews, review)
        }

        // Query: Sizes
        querySizes := "SELECT size FROM product_sizes WHERE product_id = $1"
        sizeRows, errSizes := db.GetDB().QueryContext(ctx, querySizes, productId)
        if errSizes != nil {
            fmt.Println(errSizes)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve sizes"})
            return
        }
        defer sizeRows.Close()

        for sizeRows.Next() {
            var size string
            if err := sizeRows.Scan(&size); err != nil {
                fmt.Println(err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning sizes"})
                return
            }
            sizes = append(sizes, size)
        }

        // Query: Images
        queryImages := "SELECT image_url FROM product_images WHERE product_id = $1"
        imageRows, errImages := db.GetDB().QueryContext(ctx, queryImages, productId)
        if errImages != nil {
            fmt.Println(errImages)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve images"})
            return
        }
        defer imageRows.Close()

        for imageRows.Next() {
            var image string
            if err := imageRows.Scan(&image); err != nil {
                fmt.Println(err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning images"})
                return
            }
            images = append(images, image)
        }

        // Combine all data into the product page struct
        productPage.Product = product
        productPage.Colors = colors
        productPage.Reviews = reviews
        productPage.Sizes = sizes
        productPage.Images = images

        // Respond with the full product page details
        c.JSON(http.StatusOK, productPage)
    }
}
