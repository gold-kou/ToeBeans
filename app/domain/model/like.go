package model

import "time"

type Like struct {
	ID        uint64
	UserName  string
	PostingID uint64
	CreatedAt time.Time
	UpdatedAt time.Time
}
