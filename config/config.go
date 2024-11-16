package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string
	TimeZone   string
}

type Config struct {
	DB        DBConfig
	SecretKey string
	GRPCUrl   string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &Config{
		SecretKey: os.Getenv("JWT_SECRET_KEY"),
		DB: DBConfig{
			DBHost:     os.Getenv("DB_HOST"),
			DBPort:     os.Getenv("DB_PORT"),
			DBUser:     os.Getenv("DB_USER"),
			DBPassword: os.Getenv("DB_PASSWORD"),
			DBName:     os.Getenv("DB_NAME"),
			SSLMode:    os.Getenv("DB_SSL_MODE"),
			TimeZone:   os.Getenv("DB_TIMEZONE"),
		},
	}
}
