package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"deepseen-backend/configuration"
	DB "deepseen-backend/database"
	Schemas "deepseen-backend/database/schemas"
	"deepseen-backend/utilities"
)

// Validate recovery code
func validateRecoveryCode(ctx *fiber.Ctx) error {
	// check data
	var body RecoveryValidate
	bodyParsingError := ctx.BodyParser(&body)
	if bodyParsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	code := body.Code
	password := body.Password
	if code == "" || password == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}
	trimmedCode := strings.TrimSpace(code)
	trimmedPassword := strings.TrimSpace(password)
	if trimmedCode == "" || trimmedPassword == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	// find Password
	PasswordCollection := DB.Instance.Database.Collection(DB.Collections.Password)
	rawPasswordRecord := PasswordCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "recoveryCode", Value: trimmedCode}},
	)
	passwordRecord := &Schemas.Password{}
	rawPasswordRecord.Decode(passwordRecord)
	if passwordRecord.ID == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InvalidCode,
			Status: fiber.StatusUnauthorized,
		})
	}

	// update Password record
	hash, hashError := utilities.MakeHash(trimmedPassword)
	if hashError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	now := utilities.MakeTimestamp()
	passwordID, conversionError := primitive.ObjectIDFromHex(passwordRecord.ID)
	if conversionError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	_, updateError := PasswordCollection.UpdateOne(
		ctx.Context(),
		bson.D{{Key: "_id", Value: passwordID}},
		bson.D{{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "hash",
					Value: hash,
				},
				{
					Key:   "recoveryCode",
					Value: "",
				},
				{
					Key:   "updated",
					Value: now,
				},
			},
		}},
	)
	if updateError != nil {
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
