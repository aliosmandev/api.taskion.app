package blocks

import (
	"github.com/gofiber/fiber/v2"
)

func InitRouter(router fiber.Router) {
	router.Get("/:pageId", getBlocks)
	router.Post("/addTodo/:pageId", AddTodo)
}
