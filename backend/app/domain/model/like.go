package model

import "time"

type Like struct {
	ID        int64
	UserID    int64
	PostingID int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
