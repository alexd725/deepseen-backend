package configuration

import "time"

// Clients contains available clients
var Clients = ClientsStruct{
	Desktop: "desktop",
	Mobile:  "mobile",
	Web:     "web",
}

// Environments contains available environments
var Environments = EnvironmentsStruct{
	Development: "development",
	Heroku:      "heroku",
	Production:  "production",
}

// Redis contains additional Reids configuration
var Redis = RedisOptions{
	Prefixes: RedisPrefixes{
		Room: "room",
		User: "user",
	},
	TTL: 24 * time.Hour,
}

// ResponseMessages contains server response messages
var ResponseMessages = ResponseMessagesStruct{
	AccessDenied:           "ACCESS_DENIED",
	EmailAlreadyInUse:      "EMAIL_ALREADY_IN_USE",
	InternalServerError:    "INTERNAL_SERVER_ERROR",
	InvalidCode:            "INVALID_CODE",
	InvalidData:            "INVALID_DATA",
	InvalidEmail:           "INVALID_EMAIL",
	InvalidToken:           "INVALID_TOKEN",
	MissingData:            "MISSING_DATA",
	MissingSecret:          "MISSING_SECRET",
	MissingToken:           "MISSING_TOKEN",
	MissingUserID:          "MISSING_USER_ID",
	NotFound:               "NOT_FOUND",
	Ok:                     "OK",
	OldPasswordIsInvalid:   "OLD_PASSWORD_IS_INVALID",
	PasswordRecordNotFound: "PASSWORD_RECORD_NOT_FOUND",
	TooManyRequests:        "TOO_MANY_REQUESTS",
}

// Roles contains available account roles
var Roles = RolesStruct{
	Admin: "admin",
	Root:  "root",
	User:  "user",
}
