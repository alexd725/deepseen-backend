package schemas

// Track schema structure
type Track struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`

	TrackName     string `json:"trackName"`
	TrackDuration string `json:"trackDuration"`
	UserId        string `json:"userId"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}
