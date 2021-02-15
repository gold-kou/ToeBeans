package model

import "time"

type Comment struct {
	ID        int64
	UserName  string
	PostingID int64
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
