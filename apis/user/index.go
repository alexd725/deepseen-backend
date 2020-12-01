package user

import (
	"github.com/gofiber/fiber/v2"

	"deepseen-backend/middlewares"
)

// APIs setup
func Setup(app *fiber.App) {
	group := app.Group("/api/user")

	group.Patch("/name", middlewares.Authorize, changeName)
	group.Patch("/password", middlewares.Authorize, changePassword)
}
