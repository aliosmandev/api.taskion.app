package main

import (
	"fmt"
	"go-rest-api/cmd/application"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	application.Start()
}
