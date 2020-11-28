package utilities

import (
	JWT "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type GenerateJWTParams struct {
	Client    string
	ExpiresIn int64
	Image     string
	UserId    string
}

type JWTClaims struct {
	Client string `json:"client"`
	Image  string `json:"image"`
	UserId string `json:"userId"`
	JWT.StandardClaims
}

type ResponseParams struct {
	Ctx    *fiber.Ctx
	Data   interface{}
	Info   string
	Status int
}
