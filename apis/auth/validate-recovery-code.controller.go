package auth

import (
	"github.com/gofiber/fiber/v2"

	// . "deepseen-backend/database"
	// . "deepseen-backend/database/schemas"
	"deepseen-backend/utilities"
)

// Validate recovery code
func validateRecoveryCode(ctx *fiber.Ctx) error {
	return utilities.Response(utilities.ResponseParams{
		Ctx: ctx,
	})
}
