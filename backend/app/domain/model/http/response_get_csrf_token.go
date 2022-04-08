package http

type ResponseGetCsrfToken struct {
	CsrfToken string `json:"csrf_token"`
}
