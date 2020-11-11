package model

import "time"

const (
	LikeAction    = "like"
	CommentAction = "comment"
	FollowAction  = "follow"
)

type Notification struct {
	ID          int64
	VisitorName string
	VisitedName string
	Action      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
