package user

type ChangePasswordRequest struct {
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}
