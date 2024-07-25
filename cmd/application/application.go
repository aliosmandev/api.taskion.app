package application

import (
	"github.com/gofiber/fiber/v2"

	authRoutes "taskmanager/modules/auth"
	pagesRoutes "taskmanager/modules/pages"
)

func Start() {
	app := fiber.New()
	authGroup := app.Group("/auth")
	authRoutes.InitRouter(authGroup)

	pagesGroup := app.Group("/pages")
	pagesRoutes.InitRouter(pagesGroup)

	app.Listen(":8080")
}
