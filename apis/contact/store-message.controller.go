package contact

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"deepseen-backend/configuration"
	. "deepseen-backend/database"
	. "deepseen-backend/database/schemas"
	"deepseen-backend/utilities"
)

// Store a message from the Contact form
func storeMessage(ctx *fiber.Ctx) error {
	// check data
	var body PostMessageRequest
	bodyParsingError := ctx.BodyParser(&body)
	if bodyParsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	email := body.Email
	message := body.Message
	name := body.Name
	if email == "" || message == "" || name == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}
	trimmedEmail := strings.TrimSpace(email)
	trimmedMessage := strings.TrimSpace(message)
	trimmedName := strings.TrimSpace(name)
	if trimmedEmail == "" || trimmedMessage == "" || trimmedName == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	// make sure that email address is valid
	emailIsValid := utilities.ValidateEmail(trimmedEmail)
	if !emailIsValid {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InvalidEmail,
			Status: fiber.StatusBadRequest,
		})
	}

	// create a new Message record and insert it
	MessageCollection := Instance.Database.Collection(Collections.Message)
	now := utilities.MakeTimestamp()
	NewMessage := new(Message)
	NewMessage.ID = ""
	NewMessage.Email = trimmedEmail
	NewMessage.Message = trimmedMessage
	NewMessage.Name = trimmedName
	NewMessage.Created = now
	NewMessage.Updated = now
	_, insertionError := MessageCollection.InsertOne(ctx.Context(), NewMessage)
	if insertionError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	return utilities.Response(utilities.ResponseParams{
		Ctx: ctx,
	})
}
