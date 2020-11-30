package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"

	"deepseen-backend/configuration"
	"deepseen-backend/utilities"
)

// Authorize requests for the Services APIs
func AuthorizeServices(ctx *fiber.Ctx) error {
	// get authorization header
	rawSecret := ctx.Get("X-WS-SECRET")
	if rawSecret == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingSecret,
			Status: fiber.StatusUnauthorized,
		})
	}
	trimmedSecret := strings.TrimSpace(rawSecret)
	if trimmedSecret == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingSecret,
			Status: fiber.StatusUnauthorized,
		})
	}

	// validate the secret
	wsSecret := os.Getenv("WS_SECRET")
	if wsSecret != trimmedSecret {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	return ctx.Next()
}
