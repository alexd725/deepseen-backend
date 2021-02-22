package configuration

import (
	"time"
)

// ClientsStruct represents the client types
type ClientsStruct struct {
	Desktop string
	Mobile  string
	Web     string
}

// EnvironmentsStruct represents the availabe application environments
type EnvironmentsStruct struct {
	Development string
	Heroku      string
	Production  string
}

// RedisPrefixes contains the prefixes used in the Redis keys
type RedisPrefixes struct {
	Room string
	User string
}

// RedisOptions describes the options struct
type RedisOptions struct {
	Prefixes RedisPrefixes
	TTL      time.Duration
}

// ResponseMessagesStruct describes available response messages
type ResponseMessagesStruct struct {
	AccessDenied           string
	EmailAlreadyInUse      string
	InternalServerError    string
	InvalidCode            string
	InvalidEmail           string
	InvalidData            string
	InvalidToken           string
	MissingData            string
	MissingSecret          string
	MissingToken           string
	MissingUserID          string
	NotFound               string
	Ok                     string
	OldPasswordIsInvalid   string
	PasswordRecordNotFound string
	TooManyRequests        string
}

// RolesStruct describes available roles
type RolesStruct struct {
	Admin string
	Root  string
	User  string
}
