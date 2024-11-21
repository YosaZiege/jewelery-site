package helper

import (
	"context"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/yosaZiege/jewelery-website/config"
	"github.com/yosaZiege/jewelery-website/db"
)




type SignedDetails struct{
	Email    string `json:"email"`
	Role     string `json:"role"`
	Username string `json:"username"`
	Uid      int `json:"uid"`
	jwt.StandardClaims

}

func GenerateAllTokens(username string , role string, email string , uid int) (signedToken string,refreshToken string,err error) {
	claims := &SignedDetails{
		Email: email,
		Role: role,
		Username: username,
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	signedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.GetEnv("JWT_SECRET", "")))
	if err != nil {
		log.Panic("Error generating main token:", err)
		return "", "", err
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(config.GetEnv("JWT_SECRET", "")))
		if err != nil {
			log.Panic("Error generating refresh token:", err)
			return "", "", err
		}
	return signedToken , refreshToken , nil
}
func UpdateAllTokens(signedToken string, signedRefreshToken string, uid int) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	// Parse the current time to be used for `updated_at`
	updatedAt := time.Now().Format(time.RFC3339)
	// Define the SQL query for updating the tokens and timestamp
	query := `UPDATE public.users SET token=$1, refresh_token=$2, updated_at=$3 WHERE id=$4`
	// Execute the query with the provided parameters
_, err := db.GetDB().ExecContext(ctx, query, signedToken, signedRefreshToken, updatedAt, uid)
	if err != nil {
		log.Println("Error updating tokens:", err)
		return
	} 
	log.Println("Tokens updated successfully for user ID:", uid)
}
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token , err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetEnv("JWT_SECRET", "")) , nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	} 	
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "The Token is Invalid!"
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix(){
		msg = "Token is Expired"
		msg = err.Error()
		return
	}
	return claims , msg
}