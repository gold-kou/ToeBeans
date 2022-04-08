package http

type ResponseGetNotifications struct {
	VisitedName string                    `json:"visited_name,omitempty"`
	Actions     []ResponseGetNotification `json:"actions,omitempty"`
}
