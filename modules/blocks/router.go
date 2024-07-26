package blocks

import (
	"github.com/gofiber/fiber/v2"
)

func InitRouter(router fiber.Router) {
	router.Get("/get/:pageId", getBlocks)
}
