package auth

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"deepseen-backend/configuration"
	. "deepseen-backend/database"
	. "deepseen-backend/database/schemas"
	"deepseen-backend/redis"
	"deepseen-backend/utilities"
)

// Handle complete sign out
func completeSignOut(ctx *fiber.Ctx) error {
	// get User ID from Locals (assert it as string as well)
	userId := ctx.Locals("UserId").(string)

	// load an Image record
	ImageCollection := Instance.Database.Collection(Collections.Image)
	rawImageRecord := ImageCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userId}},
	)
	imageRecord := &Image{}
	rawImageRecord.Decode(imageRecord)
	if imageRecord.ID == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	// generate a new image
	image, imageError := utilities.MakeHash(
		userId + fmt.Sprintf("%v", utilities.MakeTimestamp()),
	)
	if imageError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// update Image record
	now := utilities.MakeTimestamp()
	_, updateError := ImageCollection.UpdateOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userId}},
		bson.D{{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "image",
					Value: image,
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

	// delete room record from Redis
	redisRoomError := redis.Client.Del(
		context.Background(),
		utilities.KeyFormatter(
			configuration.Redis.Prefixes.Room,
			userId,
		),
	).Err()
	if redisRoomError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// store updated user image in Redis
	redisUserError := redis.Client.Set(
		context.Background(),
		utilities.KeyFormatter(
			configuration.Redis.Prefixes.User,
			userId,
		),
		image,
		configuration.Redis.TTL,
	).Err()
	if redisUserError != nil {
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
