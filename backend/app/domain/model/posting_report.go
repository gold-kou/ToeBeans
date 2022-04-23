package model

import "time"

type PostingReport struct {
	ID        int
	PostingID int
	Detail    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
