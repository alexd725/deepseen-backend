package user

type ChangeNameRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type ChangePasswordRequest struct {
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}
