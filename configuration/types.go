package configuration

import (
	"time"
)

type ClientsStruct struct {
	Desktop string
	Mobile  string
	Web     string
}

type EnvironmentsStruct struct {
	Development string
	Heroku      string
	Production  string
}

type RedisPrefixes struct {
	Room string
	User string
}

type RedisOptions struct {
	Prefixes RedisPrefixes
	TTL      time.Duration
}

type ResponseMessagesStruct struct {
	AccessDenied        string
	EmailAlreadyInUse   string
	InternalServerError string
	InvalidCode         string
	InvalidEmail        string
	InvalidData         string
	InvalidToken        string
	MissingData         string
	MissingSecret       string
	MissingToken        string
	MissingUserID       string
	NotFound            string
	Ok                  string
	TooManyRequests     string
}
