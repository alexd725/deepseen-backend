package manage

type ChangeRoleRequest struct {
	Role   string `json:"role"`
	UserId string `json:"userId"`
}
