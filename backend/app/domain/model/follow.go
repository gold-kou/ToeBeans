package model

import "time"

type Follow struct {
	ID              int64
	FollowingUserID int64
	FollowedUserID  int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
