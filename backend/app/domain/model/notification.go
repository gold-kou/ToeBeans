package model

import "time"

const (
	LikeAction    = "like"
	CommentAction = "comment"
	FollowAction  = "follow"
)

type Notification struct {
	ID            int64
	VisitorUserID int64
	VisitedUserID int64
	Action        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
