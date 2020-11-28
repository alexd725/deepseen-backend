package schemas

// Password schema structure
type Password struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`

	Hash         string `json:"hash"`
	RecoveryCode string `json:"recoveryCode"`
	UserId       string `json:"userId"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}