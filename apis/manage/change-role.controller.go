package manage

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"deepseen-backend/configuration"
	. "deepseen-backend/database"
	. "deepseen-backend/database/schemas"
	"deepseen-backend/utilities"
)

// Change role of a user
func changeRole(ctx *fiber.Ctx) error {
	// check data
	id := ctx.Params("id")
	if id == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingUserID,
			Status: fiber.StatusBadRequest,
		})
	}

	// parse ID into an ObjectID
	parsedId, parsingError := primitive.ObjectIDFromHex(id)
	if parsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InvalidData,
			Status: fiber.StatusNotFound,
		})
	}

	// get User record
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

	// get Image record
	ImageCollection := Instance.Database.Collection(Collections.Image)
	rawImageRecord := ImageCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: id}},
	)
	imageRecord := &Image{}
	rawImageRecord.Decode(imageRecord)
	if imageRecord.ID == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.NotFound,
			Status: fiber.StatusNotFound,
		})
	}

	return utilities.Response(utilities.ResponseParams{
		Ctx: ctx,
		Data: fiber.Map{
			"image": imageRecord,
		},
	})
}
