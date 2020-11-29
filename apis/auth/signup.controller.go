package auth

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"deepseen-backend/configuration"
	. "deepseen-backend/database"
	. "deepseen-backend/database/schemas"
	"deepseen-backend/redis"
	"deepseen-backend/utilities"
)

// Handle signing up
func signUp(ctx *fiber.Ctx) error {
	// check data
	var body SignUpUserRequest
	bodyParsingError := ctx.BodyParser(&body)
	if bodyParsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	client := body.Client
	email := body.Email
	name := body.Name
	password := body.Password
	if client == "" || email == "" || name == "" || password == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}
	trimmedClient := strings.TrimSpace(client)
	trimmedEmail := strings.TrimSpace(email)
	trimmedName := strings.TrimSpace(name)
	trimmedPassword := strings.TrimSpace(password)
	if trimmedClient == "" || trimmedEmail == "" ||
		trimmedName == "" || trimmedPassword == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	// make sure that the client is valid
	clients := utilities.Values(configuration.Clients)
	if !utilities.IncludesString(clients, trimmedClient) {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InvalidData,
			Status: fiber.StatusBadRequest,
		})
	}

	// load User schema
	UserCollection := Instance.Database.Collection(Collections.User)

	// check if email is already in use
	existingRecord := UserCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "email", Value: trimmedEmail}},
	)
	existingUser := &User{}
	existingRecord.Decode(existingUser)
	if existingUser.ID != "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.EmailAlreadyInUse,
			Status: fiber.StatusBadRequest,
		})
	}

	// create a new User record, insert it and get back the ID
	now := utilities.MakeTimestamp()
	NewUser := new(User)
	NewUser.ID = ""
	NewUser.Email = trimmedEmail
	NewUser.Name = trimmedName
	NewUser.Role = "user"
	NewUser.Created = now
	NewUser.Updated = now
	insertionResult, insertionError := UserCollection.InsertOne(ctx.Context(), NewUser)
	if insertionError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	createdRecord := UserCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "_id", Value: insertionResult.InsertedID}},
	)
	createdUser := &User{}
	createdRecord.Decode(createdUser)

	// load Image schema
	ImageCollection := Instance.Database.Collection(Collections.Image)

	// create an Image for the User
	image, imageError := utilities.MakeHash(
		createdUser.ID + fmt.Sprintf("%v", utilities.MakeTimestamp()),
	)
	if imageError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// create a new Image record and insert it
	NewImage := new(Image)
	NewImage.ID = ""
	NewImage.Image = image
	NewImage.UserId = createdUser.ID
	NewImage.Created = now
	NewImage.Updated = now
	_, insertionError = ImageCollection.InsertOne(ctx.Context(), NewImage)
	if insertionError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// load Password schema
	PasswordCollection := Instance.Database.Collection(Collections.Password)

	// create password hash
	hash, hashError := utilities.MakeHash(trimmedPassword)
	if hashError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// create a new Password record and insert it
	NewPassword := new(Password)
	NewPassword.ID = ""
	NewPassword.Hash = hash
	NewPassword.RecoveryCode = ""
	NewPassword.UserId = createdUser.ID
	NewPassword.Created = now
	NewPassword.Updated = now
	_, insertionError = PasswordCollection.InsertOne(ctx.Context(), NewPassword)
	if insertionError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// generate a token
	expiration, expirationError := strconv.Atoi(os.Getenv("TOKEN_EXPIRATION"))
	if expirationError != nil {
		expiration = 9999
	}
	token, tokenError := utilities.GenerateJWT(utilities.GenerateJWTParams{
		Client:    trimmedClient,
		ExpiresIn: int64(expiration),
		Image:     image,
		UserId:    createdUser.ID,
	})
	if tokenError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// store user image in Redis
	redisError := redis.Client.Set(
		context.Background(),
		utilities.KeyFormatter(
			configuration.Redis.Prefixes.User,
			createdUser.ID,
		),
		image,
		configuration.Redis.TTL,
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
		Data: fiber.Map{
			"token": token,
			"user":  createdUser,
		},
	})
}
