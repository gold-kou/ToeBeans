package model

import "time"

type UserReport struct {
	ID        int
	UserName  string
	Detail    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
