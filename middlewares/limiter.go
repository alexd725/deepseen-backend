package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"deepseen-backend/configuration"
	"deepseen-backend/utilities"
)

// Limiter function creates a limiter middleware
func Limiter(params LimiterParams) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        params.Max,
		Expiration: time.Duration(params.Timeframe) * time.Second,
		LimitReached: func(ctx *fiber.Ctx) error {
			return utilities.Response(utilities.ResponseParams{
				Ctx:    ctx,
				Info:   configuration.ResponseMessages.TooManyRequests,
				Status: fiber.StatusTooManyRequests,
			})
		},
	})
}
