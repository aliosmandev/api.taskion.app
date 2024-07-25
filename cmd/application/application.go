package application

import (
	"github.com/gofiber/fiber/v2"

	authRoutes "taskmanager/modules/auth"
)

func Start() {
	app := fiber.New()
	authGroup := app.Group("/auth")
	authRoutes.InitRouter(authGroup)

	app.Listen(":8080")
}
