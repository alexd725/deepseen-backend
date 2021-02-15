package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
	"github.com/joho/godotenv"

	"deepseen-backend/apis/auth"
	"deepseen-backend/apis/contact"
	"deepseen-backend/apis/manage"
	"deepseen-backend/apis/services"
	"deepseen-backend/apis/user"
	"deepseen-backend/configuration"
	"deepseen-backend/database"
	"deepseen-backend/middlewares"
	"deepseen-backend/redis"
	"deepseen-backend/utilities"
)

func main() {
	// load environment variables via the .env file
	env := os.Getenv("ENV")
	if env != configuration.Environments.Heroku {
		envError := godotenv.Load()
		if envError != nil {
			log.Fatal(envError)
			return
		}
	}

	// connect to the database
	databaseError := database.Connect()
	if databaseError != nil {
		log.Fatal(databaseError)
		return
	}

	// connect to the Redis server
	redisError := redis.Connect()
	if redisError != nil {
		log.Fatal(redisError)
		return
	}

	app := fiber.New()

	// middlewares
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(favicon.New(favicon.Config{
		File: "./assets/favicon.ico",
	}))
	app.Use(helmet.New())
	app.Use(middlewares.Limiter(middlewares.LimiterParams{
		Max:       90,
		Timeframe: 60,
	}))
	app.Use(logger.New())

	// available APIs
	auth.Setup(app)
	contact.Setup(app)
	manage.Setup(app)
	services.Setup(app)
	user.Setup(app)

	// handle 404
	app.Use(func(ctx *fiber.Ctx) error {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.NotFound,
			Status: fiber.StatusNotFound,
		})
	})

	// get the port
	port := os.Getenv("PORT")
	if port == "" {
		port = "1337"
	}

	// launch the app
	launchError := app.Listen(":" + port)
	if launchError != nil {
		panic(launchError)
	}
}
