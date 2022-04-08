package http

type ResponseSimpleSuccess struct {
	Status  int32  `json:"status"`
	Message string `json:"message"`
}
