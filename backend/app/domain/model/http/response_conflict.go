package http

type ResponseConflict struct {
	Status  int32  `json:"status"`
	Message string `json:"message"`
}
