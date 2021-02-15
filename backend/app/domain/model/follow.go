package model

import "time"

type Follow struct {
	ID                int64
	FollowingUserName string
	FollowedUserName  string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
