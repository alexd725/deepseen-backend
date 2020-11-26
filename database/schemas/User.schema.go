package schemas

// User schema structure
type User struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`

	Email string `json:"email"`
	Image string `json:"image"`
	Name  string `json:"name"`
	Role  string `json:"role"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}
