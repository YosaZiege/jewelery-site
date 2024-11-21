package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(){
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Failed To load .env File")
	}
}

func GetEnv(key , defaultValue string) string{
	value , exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}