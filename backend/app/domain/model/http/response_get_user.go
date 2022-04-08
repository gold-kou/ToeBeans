package http

import (
	"time"
)

type ResponseGetUser struct {
	UserName         string    `json:"user_name"`
	Icon             string    `json:"icon"`
	SelfIntroduction string    `json:"self_introduction"`
	PostingCount     int64     `json:"posting_count"`
	LikeCount        int64     `json:"like_count"`
	LikedCount       int64     `json:"liked_count"`
	FollowCount      int64     `json:"follow_count"`
	FollowedCount    int64     `json:"followed_count"`
	CreatedAt        time.Time `json:"created_at"`
}
