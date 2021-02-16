package manage

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"deepseen-backend/configuration"
	. "deepseen-backend/database"
	"deepseen-backend/utilities"
)

// Change role of a user
func changeRole(ctx *fiber.Ctx) error {
	// check data
	var body ChangeRoleRequest
	bodyParsingError := ctx.BodyParser(&body)
	if bodyParsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	role := body.Role
	userId := body.UserId
	if role == "" || userId == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}
	trimmedRole := strings.TrimSpace(role)
	trimmedUserId := strings.TrimSpace(userId)
	if trimmedRole == "" || trimmedUserId == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	// check role validity
	roles := utilities.Values(configuration.Roles)
	if !utilities.IncludesString(roles, trimmedRole) {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InvalidData,
			Status: fiber.StatusBadRequest,
		})
	}

	// parse user ID into an ObjectID
	parsedId, parsingError := primitive.ObjectIDFromHex(userId)
	if parsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InvalidData,
			Status: fiber.StatusNotFound,
		})
	}

	// update User role
	now := utilities.MakeTimestamp()
	UserCollection := Instance.Database.Collection(Collections.User)
	_, updateError := UserCollection.UpdateOne(
		ctx.Context(),
		bson.D{{Key: "_id", Value: parsedId}},
		bson.D{{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "role",
					Value: trimmedRole,
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
