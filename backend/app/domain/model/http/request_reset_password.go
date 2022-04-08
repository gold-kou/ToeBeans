package http

type RequestResetPassword struct {
	UserName         string `json:"user_name"`
	Password         string `json:"password"`
	PasswordResetKey string `json:"password_reset_key"`
}
