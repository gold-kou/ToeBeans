package http

type ResponseInternalServerError struct {
	Status  int32  `json:"status"`
	Message string `json:"message"`
}
