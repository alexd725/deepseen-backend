package auth

type RecoveryEmail struct {
	Email string `json:"email"`
}

type RecoveryValidate struct {
	Code     string `json:"code"`
	Password string `json:"password"`
}

type SignInUserRequest struct {
	Client   string `json:"client"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	SignInUserRequest
}
