package configuration

import "time"

// Available clients
var Clients = ClientsStruct{
	Desktop: "desktop",
	Mobile:  "mobile",
	Web:     "web",
}

// Additional Reids configuration
var Redis = RedisOptions{
	Prefixes: RedisPrefixes{
		Room: "room",
		User: "user",
	},
	TTL: 24 * time.Hour,
}

// Server response messages
var ResponseMessages = ResponseMessagesStruct{
	AccessDenied:        "ACCESS_DENIED",
	EmailAlreadyInUse:   "EMAIL_ALREADY_IN_USE",
	InternalServerError: "INTERNAL_SERVER_ERROR",
	InvalidCode:         "INVALID_CODE",
	InvalidData:         "INVALID_DATA",
	InvalidToken:        "INVALID_TOKEN",
	MissingData:         "MISSING_DATA",
	MissingSecret:       "MISSING_SECRET",
	MissingToken:        "MISSING_TOKEN",
	MissingUserID:       "MISSING_USER_ID",
	NotFound:            "NOT_FOUND",
	Ok:                  "OK",
	TooManyRequests:     "TOO_MANY_REQUESTS",
}
