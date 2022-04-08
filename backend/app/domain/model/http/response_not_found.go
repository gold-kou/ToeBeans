package http

type ResponseNotFound struct {
	Status  int32  `json:"status"`
	Message string `json:"message"`
}
