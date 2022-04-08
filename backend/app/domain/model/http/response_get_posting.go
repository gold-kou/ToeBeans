package http

import (
	"time"
)

type ResponseGetPosting struct {
	PostingId  int64     `json:"posting_id"`
	UserName   string    `json:"user_name"`
	UploadedAt time.Time `json:"uploaded_at"`
	Title      string    `json:"title"`
	ImageUrl   string    `json:"image_url,omitempty"`
	LikedCount int64     `json:"liked_count"`
	Liked      bool      `json:"liked"`
}
