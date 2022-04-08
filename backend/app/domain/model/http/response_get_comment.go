package http

import (
	"time"
)

type ResponseGetComment struct {
	CommentId   int64     `json:"comment_id"`
	UserName    string    `json:"user_name"`
	CommentedAt time.Time `json:"commented_at"`
	Comment     string    `json:"comment"`
}
