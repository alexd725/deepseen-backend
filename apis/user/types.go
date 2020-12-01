package user

type ChangeNameRequest struct {
	Name string `json:"name"`
}

type ChangePasswordRequest struct {
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}
