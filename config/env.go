package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func MongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("MONGOURI")
}
