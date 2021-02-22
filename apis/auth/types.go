package auth

// RecoveryEmail describes a body of the Send request
type RecoveryEmail struct {
	Email string `json:"email"`
}

// RecoveryValidate describes a body of the Validate request
type RecoveryValidate struct {
	Code     string `json:"code"`
	Password string `json:"password"`
}

// SignInUserRequest describes a body of the Sign In request
type SignInUserRequest struct {
	Client   string `json:"client"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUpUserRequest describes a body of the Sign Up request
type SignUpUserRequest struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	SignedAgreement bool   `json:"signedAgreement"`
	SignInUserRequest
}
