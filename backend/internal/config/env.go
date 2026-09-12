package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Loads .env when available and reports errors for malformed files.
func LoadEnv() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("config: failed to load .env: %v", err)
	}
}
