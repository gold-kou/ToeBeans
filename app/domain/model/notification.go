package model

import "time"

type Notification struct {
	ID          uint64
	VisitorName string
	VisitedName string
	PostingID   uint64
	CommentID   uint64
	Action      string
	Checked     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
