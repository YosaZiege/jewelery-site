package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	helper "github.com/yosaZiege/jewelery-website/helpers"
)






func Authenticate() gin.HandlerFunc{
	return func (c *gin.Context)  {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"Error" : "No Token FOund"})
			c.Abort()
			return
		}

		claims , err := helper.ValidateToken(clientToken)

		if err != ""{
			c.JSON(http.StatusInternalServerError , gin.H{"Error" : "Token is Invalid"})
			c.Abort()
			return
		}

		c.Set("email" , claims.Email)
		c.Set("role" , claims.Role)
		c.Set("username" , claims.Username)
		c.Set("uid" , claims.Uid)
		c.Next()
	}
}