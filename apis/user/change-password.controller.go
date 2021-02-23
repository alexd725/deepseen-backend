package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"deepseen-backend/configuration"
	. "deepseen-backend/database"
	. "deepseen-backend/database/schemas"
	"deepseen-backend/redis"
	"deepseen-backend/utilities"
)

// Change a password
func changePassword(ctx *fiber.Ctx) error {
	// check data
	var body ChangePasswordRequest
	bodyParsingError := ctx.BodyParser(&body)
	if bodyParsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	newPassword := body.NewPassword
	oldPassword := body.OldPassword
	if newPassword == "" || oldPassword == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}
	trimmedNewPassword := strings.TrimSpace(newPassword)
	trimmedOldPassword := strings.TrimSpace(oldPassword)
	if trimmedNewPassword == "" || trimmedOldPassword == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	// get User ID from Locals (assert it as string as well)
	userId := ctx.Locals("UserId").(string)

	// load Password record
	PasswordCollection := Instance.Database.Collection(Collections.Password)
	rawPasswordRecord := PasswordCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userId}},
	)
	passwordRecord := &Password{}
	rawPasswordRecord.Decode(passwordRecord)
	if passwordRecord.ID == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.PasswordRecordNotFound,
			Status: fiber.StatusUnauthorized,
		})
	}

	// compare hashes
	passwordIsValid := utilities.CompareHashes(trimmedOldPassword, passwordRecord.Hash)
	if !passwordIsValid {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.OldPasswordIsInvalid,
			Status: fiber.StatusUnauthorized,
		})
	}

	// generate a new password hash
	hash, hashError := utilities.MakeHash(newPassword)
	if hashError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// update Password record
	now := utilities.MakeTimestamp()
	_, updateError := PasswordCollection.UpdateOne(
		ctx.Context(),
		bson.D{{Key: "_id", Value: passwordRecord.ID}},
		bson.D{{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "hash",
					Value: hash,
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

	// update Image record as well
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
			Info:   configuration.ResponseMessages.ImageRecordNotFound,
			Status: fiber.StatusUnauthorized,
		})
	}
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
	_, updateError = ImageCollection.UpdateOne(
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

	// delete image record from Redis
	redisUserError := redis.Client.Del(
		context.Background(),
		utilities.KeyFormatter(
			configuration.Redis.Prefixes.User,
			userId,
		),
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
