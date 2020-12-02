package user

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"deepseen-backend/configuration"
	. "deepseen-backend/database"
	. "deepseen-backend/database/schemas"
	"deepseen-backend/utilities"
)

// Get user record
func getUser(ctx *fiber.Ctx) error {
	// get User ID from Locals (assert it as string as well)
	userId := ctx.Locals("UserId").(string)

	// parse ID into an ObjectID
	parsedId, parsingError := primitive.ObjectIDFromHex(userId)
	if parsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// load record
	UserCollection := Instance.Database.Collection(Collections.User)
	rawUserRecord := UserCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "_id", Value: parsedId}},
	)
	userRecord := &User{}
	rawUserRecord.Decode(userRecord)
	if userRecord.ID == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.NotFound,
			Status: fiber.StatusNotFound,
		})
	}

	return utilities.Response(utilities.ResponseParams{
		Ctx: ctx,
		Data: fiber.Map{
			"user": userRecord,
		},
	})
}
