package schemas

// Message schema structure
type Message struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`

	Email   string `json:"email"`
	Message string `json:"message"`
	Name    string `json:"name"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}
