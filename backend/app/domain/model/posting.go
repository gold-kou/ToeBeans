package model

import "time"

type Posting struct {
	ID        int64
	UserName  string
	Title     string
	ImageURL  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
