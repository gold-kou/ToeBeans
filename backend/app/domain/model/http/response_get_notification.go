package http

import (
	"time"
)

type ResponseGetNotification struct {
	VisitorName string    `json:"visitor_name"`
	PostingId   int64     `json:"posting_id,omitempty"`
	CommentId   int64     `json:"comment_id,omitempty"`
	ActionType  string    `json:"action_type,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
