package configuration

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

type ClientsStruct struct {
	Desktop string
	Mobile  string
	Web     string
}
