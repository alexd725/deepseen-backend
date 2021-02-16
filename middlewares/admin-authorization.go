package middlewares

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

// AuthorizeAdmin function authorizes admin requests for the Manage APIs
func AuthorizeAdmin(ctx *fiber.Ctx) error {
	// get authorization header
	rawToken := ctx.Get("Authorization")
	if rawToken == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingToken,
			Status: fiber.StatusUnauthorized,
		})
	}
	trimmedToken := strings.TrimSpace(rawToken)
	if trimmedToken == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingToken,
			Status: fiber.StatusUnauthorized,
		})
	}

	// parse JWT
	claims, parsingError := utilities.ParseClaims(trimmedToken)
	if parsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	// load an Image record
	ImageCollection := DB.Instance.Database.Collection(DB.Collections.Image)
	rawImageRecord := ImageCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: claims.UserId}},
	)
	imageRecord := &Schemas.Image{}
	rawImageRecord.Decode(imageRecord)
	if imageRecord.ID == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	// compare images
	if claims.Image != imageRecord.Image {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	// parse ID into an ObjectID
	parsedID, parsingError := primitive.ObjectIDFromHex(claims.UserId)
	if parsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// load User record
	UserCollection := DB.Instance.Database.Collection(DB.Collections.User)
	rawUserRecord := UserCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "_id", Value: parsedID}},
	)
	userRecord := &Schemas.User{}
	rawUserRecord.Decode(userRecord)
	if userRecord.ID == "" ||
		!(userRecord.Role == configuration.Roles.Admin ||
			userRecord.Role == configuration.Roles.Root) {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	// store client and token data in Locals
	ctx.Locals("Client", claims.Client)
	ctx.Locals("UserId", claims.UserId)
	return ctx.Next()
}
