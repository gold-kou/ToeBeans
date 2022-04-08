package http

type ResponseNotAllowedMethod struct {
	Status  int32  `json:"status"`
	Message string `json:"message"`
}
