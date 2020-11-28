package schemas

// Track schema structure
type Track struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`

	TrackName     string `json:"trackName" bson:"trackName"`
	TrackDuration string `json:"trackDuration" bson:"trackDuration"`
	UserId        string `json:"userId" bson:"userId"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}
