package manage

import (
	"github.com/gofiber/fiber/v2"
	// "deepseen-backend/middlewares"
)

// APIs setup
func Setup(app *fiber.App) {
	group := app.Group("/api/manage")

	group.Post("/role", changeRole)
}
