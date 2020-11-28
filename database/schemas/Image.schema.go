package schemas

// Image schema structure
type Image struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`

	Image  string `json:"image"`
	UserId string `json:"userId" bson:"userId"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}
