package user

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"deepseen-backend/configuration"
	. "deepseen-backend/database"
	"deepseen-backend/redis"
	"deepseen-backend/utilities"
)

// Delete account
func deleteAccount(ctx *fiber.Ctx) error {
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

	// delete User record
	UserCollection := Instance.Database.Collection(Collections.User)
	_, deleteError := UserCollection.DeleteOne(
		ctx.Context(),
		bson.D{{Key: "_id", Value: parsedId}},
	)
	if deleteError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// delete Password record
	PasswordCollection := Instance.Database.Collection(Collections.Password)
	_, deleteError = PasswordCollection.DeleteOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userId}},
	)
	if deleteError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// delete Image record
	ImageCollection := Instance.Database.Collection(Collections.Image)
	_, deleteError = ImageCollection.DeleteOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userId}},
	)
	if deleteError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// delete room record from Redis
	redisError := redis.Client.Del(
		context.Background(),
		utilities.KeyFormatter(
			configuration.Redis.Prefixes.Room,
			userId,
		),
	).Err()
	if redisError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// delete Image record from Redis
	redisError = redis.Client.Del(
		context.Background(),
		utilities.KeyFormatter(
			configuration.Redis.Prefixes.User,
			userId,
		),
	).Err()
	if redisError != nil {
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
