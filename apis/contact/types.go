package contact

// PostMessageRequest describes the body of the Contact POST request
type PostMessageRequest struct {
	Email   string `json:"email"`
	Message string `json:"message"`
	Name    string `json:"name"`
}
