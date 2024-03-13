package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ENV struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string

	Table string
	Redis string
}

var Env ENV

func GetEnv() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Read database configuration from .env file
	Env.Host = os.Getenv("DB_HOST")
	Env.Port, _ = strconv.Atoi(os.Getenv("DB_PORT")) // Convert port to int
	Env.User = os.Getenv("DB_USER")
	Env.Password = os.Getenv("DB_PASSWORD")
	Env.DBName = os.Getenv("DB_NAME")

	Env.Redis = os.Getenv("REDIS_URL")
	Env.Table = os.Getenv("TABLE_URL")
}
