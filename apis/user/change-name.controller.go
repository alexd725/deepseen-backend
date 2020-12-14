package user

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"deepseen-backend/configuration"
	. "deepseen-backend/database"
	"deepseen-backend/utilities"
)

// Change a name
func changeName(ctx *fiber.Ctx) error {
	// check data
	var body ChangeNameRequest
	bodyParsingError := ctx.BodyParser(&body)
	if bodyParsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	firstName := body.FirstName
	lastName := body.LastName
	if firstName == "" || lastName == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}
	trimmedFirstName := strings.TrimSpace(firstName)
	trimmedLastName := strings.TrimSpace(lastName)
	if trimmedFirstName == "" || trimmedLastName == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

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

	// update User record
	UserCollection := Instance.Database.Collection(Collections.User)
	now := utilities.MakeTimestamp()
	_, updateError := UserCollection.UpdateOne(
		ctx.Context(),
		bson.D{{Key: "_id", Value: parsedId}},
		bson.D{{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "firstName",
					Value: firstName,
				},
				{
					Key:   "lastName",
					Value: lastName,
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
