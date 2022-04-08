package http

type RequestRegisterUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
