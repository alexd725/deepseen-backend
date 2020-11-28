package configuration

import (
	"time"
)

type ClientsStruct struct {
	Desktop string
	Mobile  string
	Web     string
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
	InvalidData         string
	InvalidToken        string
	MissingData         string
	MissingToken        string
	NotFound            string
	Ok                  string
	TooManyRequests     string
}
