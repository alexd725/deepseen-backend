package contact

type PostMessageRequest struct {
	Email   string `json:"email"`
	Message string `json:"message"`
	Name    string `json:"name"`
}
