package middlewares

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"deepseen-backend/configuration"
	. "deepseen-backend/database"
	. "deepseen-backend/database/schemas"
	"deepseen-backend/redis"
	"deepseen-backend/utilities"
)

// Authorize requests for the general APIs
func Authorize(ctx *fiber.Ctx) error {
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

	// check Redis
	key := utilities.KeyFormatter(
		configuration.Redis.Prefixes.User,
		claims.UserId,
	)
	redisContext := context.Background()
	redisImage, redisError := redis.Client.Get(redisContext, key).Result()
	if redisError != nil {
		// the key was not found
		if redisError == redis.Nil {
			// load an Image record
			ImageCollection := Instance.Database.Collection(Collections.Image)
			rawImageRecord := ImageCollection.FindOne(
				ctx.Context(),
				bson.D{{Key: "userId", Value: claims.UserId}},
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

			// store image in Redis regardless of its validity
			redisUserError := redis.Client.Set(
				redisContext,
				key,
				imageRecord.Image,
				configuration.Redis.TTL,
			).Err()
			if redisUserError != nil {
				return utilities.Response(utilities.ResponseParams{
					Ctx:    ctx,
					Info:   configuration.ResponseMessages.InternalServerError,
					Status: fiber.StatusInternalServerError,
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

			// store token data in Locals
			ctx.Locals("Client", claims.Client)
			ctx.Locals("UserId", claims.UserId)
			return ctx.Next()
		}
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	if redisImage != claims.Image {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	// update EXPIRE for the record in Redis
	expireError := redis.Client.Expire(redisContext, key, configuration.Redis.TTL).Err()
	if expireError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// store token data in Locals
	ctx.Locals("Client", claims.Client)
	ctx.Locals("UserId", claims.UserId)
	return ctx.Next()
}
