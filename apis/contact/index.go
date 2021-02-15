package contact

import (
	"github.com/gofiber/fiber/v2"

	"deepseen-backend/middlewares"
)

// APIs setup
func Setup(app *fiber.App) {
	group := app.Group("/api/contact")

	group.Post(
		"/",
		middlewares.Limiter(middlewares.LimiterParams{
			Max:       5,
			Timeframe: 60 * 60 * 12,
		}),
		storeMessage,
	)
}
