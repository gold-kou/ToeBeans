package model

import "time"

type Posting struct {
	ID         int64
	UserName   string
	Title      string
	ImageURL   string
	LikedCount int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
