package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv memuat variabel lingkungan dari file .env
func LoadEnv() {
	// Load hanya dalam environment development
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file, using environment variables")
		} else {
			log.Println("Environment variables loaded from .env file")
		}
	}

	// Set default values jika environment variable tidak ada
	if os.Getenv("DB_HOST") == "" {
		os.Setenv("DB_HOST", "localhost")
	}
	if os.Getenv("DB_PORT") == "" {
		os.Setenv("DB_PORT", "5432")
	}
	if os.Getenv("DB_USER") == "" {
		os.Setenv("DB_USER", "postgres")
	}
	if os.Getenv("DB_NAME") == "" {
		os.Setenv("DB_NAME", "ecommerce")
	}
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "your-secret-key")
	}
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8080")
	}
} 