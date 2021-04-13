package auth

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type restoreRequest struct {
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	Code       string `json:"code"`
}

type restoreValidRequest struct {
	Email string `json:"email"`
	Code string `json:"code"`
}

type validRequest struct {
	Token string `json:"token"`
}
