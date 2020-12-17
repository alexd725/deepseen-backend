package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lucsky/cuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"deepseen-backend/configuration"
	. "deepseen-backend/database"
	. "deepseen-backend/database/schemas"
	"deepseen-backend/utilities"
)

// Send an email with recovery link
func sendRecoveryEmail(ctx *fiber.Ctx) error {
	// check data
	var body RecoveryEmail
	bodyParsingError := ctx.BodyParser(&body)
	if bodyParsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	email := body.Email
	if email == "" || strings.TrimSpace(email) == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	// find User
	UserCollection := Instance.Database.Collection(Collections.User)
	rawUserRecord := UserCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "email", Value: strings.TrimSpace(email)}},
	)
	userRecord := &User{}
	rawUserRecord.Decode(userRecord)
	if userRecord.ID == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	// find Password
	PasswordCollection := Instance.Database.Collection(Collections.Password)
	rawPasswordRecord := PasswordCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userRecord.ID}},
	)
	passwordRecord := &Password{}
	rawPasswordRecord.Decode(passwordRecord)
	if passwordRecord.ID == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	// generate a recovery code with CUID
	recoveryCode := cuid.Slug()

	// update the Password record
	now := utilities.MakeTimestamp()
	passwordId, conversionError := primitive.ObjectIDFromHex(passwordRecord.ID)
	if conversionError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	_, updateError := PasswordCollection.UpdateOne(
		ctx.Context(),
		bson.D{{Key: "_id", Value: passwordId}},
		bson.D{{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "recoveryCode",
					Value: recoveryCode,
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

	// send an email with recovery link
	formattedTemplate := utilities.CreateRecoveryTemplate(
		recoveryCode,
		userRecord.FirstName,
		userRecord.LastName,
	)
	utilities.SendEmail(userRecord, "Deepseen: account recovery", formattedTemplate)

	return utilities.Response(utilities.ResponseParams{
		Ctx: ctx,
	})
}
