package model

import "time"

type Like struct {
	ID        int64
	UserName  string
	PostingID int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
