package main

import (
	"log"
	"os"
	"path/filepath"
	"taskmanager/cmd/application"

	"github.com/joho/godotenv"
)

func main() {
	LoadENV()
	application.Start()
}

func LoadENV() {
	// Ortam değişkenini oku
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development" // Default olarak development
	}

	var envPath string
	var err error

	if env == "production" {
		executable, err := os.Executable()
		if err != nil {
			log.Fatalf("Error getting executable path: %s", err)
		}
		exPath := filepath.Join(filepath.Dir(executable), "../api")
		err = os.Chdir(exPath)
		if err != nil {
			log.Fatalf("Error changing working directory: %s", err)
		}
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current working directory: %s", err)
		}
		envPath = filepath.Join(cwd, ".env")
	} else {
		envPath = ".env"
	}

	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file from %s: %s", envPath, err)
	}
	log.Printf("Loaded .env file from %s", envPath)
}
