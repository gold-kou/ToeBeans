package http

type RequestUpdateUser struct {
	Password         string `json:"password,omitempty"`
	Icon             string `json:"icon,omitempty"`
	SelfIntroduction string `json:"self_introduction,omitempty"`
}
