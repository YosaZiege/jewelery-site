package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"github.com/yosaZiege/jewelery-website/db"
	helper "github.com/yosaZiege/jewelery-website/helpers"
	"github.com/yosaZiege/jewelery-website/models"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()



func HashPassword(password string) string{
	bytes ,err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
    
}
func VerifyPassword(userPassword string,providedPassword string)(bool,string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword),[]byte(providedPassword))
	check := true
	msg := ""
	if err != nil {
		fmt.Println(err)
		msg = "Email or Password Incorrect"
		check = false
	}
	return check,msg
}
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		// Bind JSON input to the user struct
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the input data
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		password := HashPassword(user.Password)
		user.Password = password
		// Insert user into the database and retrieve the auto-generated ID
		err := db.GetDB().QueryRowContext(ctx, 
			"INSERT INTO users (username, password, email, role) VALUES ($1, $2, $3, $4) RETURNING id",
			user.Username, user.Password, user.Email, user.Role).Scan(&user.ID)

		if err != nil {
			// Check for unique constraint violations or other errors
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{"error": "User with this email or username already exists"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register user"})
			}
			return
		}

		// Generate tokens using the retrieved ID
		token, refreshToken, _ := helper.GenerateAllTokens(user.Username, user.Role, user.Email, user.ID)
		user.Token = token
		user.RefreshToken = refreshToken

		// Update the user with tokens
		 err = db.GetDB().QueryRowContext(ctx, 
			"UPDATE users SET token = $1, refresh_token = $2 WHERE id = $3 RETURNING token , refresh_token , role;", 
			token, refreshToken, user.ID).Scan(&user.Token , &user.RefreshToken , &user.Role)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not set tokens"})
			return
		}

		
		c.JSON(http.StatusOK, gin.H{
			"message":       "Register successful",
			"token":         user.Token,
			"refresh_token": user.RefreshToken,
			"role": user.Role,
			"userId": user.ID,
		})
	}
}



func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData struct {
			Username    string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		ctx , cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		query := "SELECT id, username, password,email,role,image_url FROM public.users where username=$1 "

		err := db.GetDB().QueryRowContext(ctx,query,loginData.Username).Scan(&user.ID,&user.Username,&user.Password,&user.Email,&user.Role,&user.ImageUrl)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"Error" : "Invalid email or password"})
			}
			log.Println("Database error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Pasword Verification Here
		passwordIsValid , msg := VerifyPassword(user.Password , loginData.Password)
		defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"Error" : msg})
			return
		}
		// Lets Generete some tokens
		token , refreshToken , err := helper.GenerateAllTokens(user.Username , user.Role, user.Email, user.ID)
		if err != nil {
			log.Println("Token generation error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not log in"})
			return
		}
		helper.UpdateAllTokens(token,refreshToken, user.ID)
		user.Token = token
		user.RefreshToken = refreshToken
		_, err = db.GetDB().ExecContext(ctx, "UPDATE public.users SET token=$1, refresh_token=$2 WHERE id=$3", user.Token, user.RefreshToken, user.ID)
		if err != nil {
			log.Println("Database error updating tokens:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not complete login"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":       "Login successful",
			"token":         token,
			"refresh_token": refreshToken,
			"role": user.Role,
			"userId": user.ID,
			"image_url": user.ImageUrl,
		})
	}
}
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var users []models.User

		// Query to select all users
		query := "SELECT id , username, email, role, is_email_verified FROM public.users;"

		// Execute the query
		rows, err := db.GetDB().QueryContext(ctx, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve users"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var user models.User

			if err := rows.Scan(&user.ID,&user.Username,&user.Email,&user.Role,&user.IsEmailVerified); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error" : "Error retrieving row"})
				return
			}
			users = append(users, user)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving users"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}


func GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		fmt.Println(userId)
		// Check user type and UID
		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return	
		}

		// Set context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		query := "SELECT username, password, email, role, image_url FROM public.users WHERE id = $1;"
		fmt.Println(query)
		// Execute query and scan result into user struct
		err := db.GetDB().QueryRowContext(ctx, query, userId).Scan(&user.Username, &user.Password, &user.Email, &user.Role, &user.ImageUrl)
		fmt.Println(err)
		if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
		}

		// Respond with user data
		c.JSON(http.StatusOK, user)
	}
}
func GetUserByEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		userEmail := c.Param("user_email")
		ctx , cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		query := "SELECT id, username , password , email , role FROM public.users where email=$1;"

		err := db.GetDB().QueryRowContext(ctx , query , userEmail).Scan(&user.ID,&user.Username , &user.Password, &user.Email,&user.Role)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusNotFound , gin.H{"error" : "User not found !"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func GetUsersV2() gin.HandlerFunc{
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c , "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page , err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}
		startIndex := (page - 1 )* recordPerPage
		query := `
			SELECT id, username, email, role, created_at, updated_at
			FROM public.users
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2`

		rows, err := db.GetDB().QueryContext(ctx, query , recordPerPage , startIndex)
		if err != nil {
			c.JSON(http.StatusInternalServerError , gin.H{"Error" : "Could not retrieve users"})
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID , &user.Username, &user.Email ,&user.Role, &user.CreatedAt , &user.UpdatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error" : "Error parsing Data"})
				return
			}
			users = append(users, user)

		}
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating over user data"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"page":         page,
			"recordCount":  len(users),
			"recordPerPage": recordPerPage,
			"data":         users,
		})
	}
}