package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key, defaultValue string) string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not loaded in processor")
	}
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
