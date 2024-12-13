package environment

import (
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	// Loading spare .env.local
	err := godotenv.Load(".env.example")
	if err != nil {
		panic("Error loading .env.example file")
	}

	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return ""
}
