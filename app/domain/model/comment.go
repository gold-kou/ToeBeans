package model

import "time"

type Comment struct {
	ID        uint64
	UserName  string
	PostingID uint64
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
