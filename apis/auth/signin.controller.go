package auth

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"deepseen-backend/configuration"
	DB "deepseen-backend/database"
	Schemas "deepseen-backend/database/schemas"
	"deepseen-backend/redis"
	"deepseen-backend/utilities"
)

// Handle signing in
func signIn(ctx *fiber.Ctx) error {
	// check data
	var body SignInUserRequest
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
	password := body.Password
	if client == "" || email == "" || password == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}
	trimmedClient := strings.TrimSpace(client)
	trimmedEmail := strings.TrimSpace(email)
	trimmedPassword := strings.TrimSpace(password)
	if trimmedClient == "" || trimmedEmail == "" || trimmedPassword == "" {
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
	UserCollection := DB.Instance.Database.Collection(DB.Collections.User)

	// find a user
	rawUserRecord := UserCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "email", Value: trimmedEmail}},
	)
	userRecord := &Schemas.User{}
	rawUserRecord.Decode(userRecord)
	if userRecord.ID == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	// load Password schema
	PasswordCollection := DB.Instance.Database.Collection(DB.Collections.Password)

	// find a password
	rawPasswordRecord := PasswordCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userRecord.ID}},
	)
	passwordRecord := &Schemas.Password{}
	rawPasswordRecord.Decode(passwordRecord)
	if passwordRecord.ID == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	// compare hashes
	passwordIsValid := utilities.CompareHashes(trimmedPassword, passwordRecord.Hash)
	if !passwordIsValid {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	// load Image schema
	ImageCollection := DB.Instance.Database.Collection(DB.Collections.Image)

	// find an image
	rawImageRecord := ImageCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userRecord.ID}},
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

	// generate a token
	expiration, expirationError := strconv.Atoi(os.Getenv("TOKEN_EXPIRATION"))
	if expirationError != nil {
		expiration = 9999
	}
	token, tokenError := utilities.GenerateJWT(utilities.GenerateJWTParams{
		Client:    trimmedClient,
		ExpiresIn: int64(expiration),
		Image:     imageRecord.Image,
		UserId:    userRecord.ID,
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
			userRecord.ID,
		),
		imageRecord.Image,
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
			"user":  userRecord,
		},
	})
}
