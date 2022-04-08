package http

type ResponseUnauthorized struct {
	Status  int32  `json:"status"`
	Message string `json:"message"`
}
