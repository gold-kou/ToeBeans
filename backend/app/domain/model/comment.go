package model

import "time"

type Comment struct {
	ID        int64
	UserID    int64
	PostingID int64
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
