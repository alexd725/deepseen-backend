package auth

import (
	"github.com/gofiber/fiber/v2"

	"deepseen-backend/middlewares"
)

// APIs setup
func Setup(app *fiber.App) {
	group := app.Group("/api/auth")

	group.Post(
		"/recovery/send",
		middlewares.Limiter(middlewares.LimiterParams{
			Max:       2,
			Timeframe: 60 * 60 * 12,
		}),
		sendRecoveryEmail,
	)

	group.Post("/recovery/validate", validateRecoveryCode)

	group.Get(
		"/signout/complete",
		middlewares.Authorize,
		completeSignOut,
	)

	group.Post(
		"/signin",
		middlewares.Limiter(middlewares.LimiterParams{
			Max:       5,
			Timeframe: 60 * 60 * 1,
		}),
		signIn,
	)

	group.Post(
		"/signup",
		middlewares.Limiter(middlewares.LimiterParams{
			Max:       1,
			Timeframe: 60 * 60 * 1,
		}),
		signUp,
	)
}
