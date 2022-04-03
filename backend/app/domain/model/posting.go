package model

import "time"

type Posting struct {
	ID        int64
	UserID    int64
	Title     string
	ImageURL  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
