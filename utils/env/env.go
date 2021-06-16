package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env.")
	}
}

func Get(key string) string {
	return os.Getenv(key)
}

