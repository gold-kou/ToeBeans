package http

type ResponseForbidden struct {
	Status  int32  `json:"status"`
	Message string `json:"message"`
}
