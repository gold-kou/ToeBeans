package model

import "time"

type User struct {
	Name                      string
	Email                     string
	Password                  string
	Icon                      string
	SelfIntroduction          string
	PostingCount              int64 // int64 because of openapi-generator
	LikeCount                 int64 // int64 because of openapi-generator
	LikedCount                int64 // int64 because of openapi-generator
	FollowCount               int64 // int64 because of openapi-generator
	FollowedCount             int64 // int64 because of openapi-generator
	ActivationKey             string
	EmailVerified             bool
	PasswordResetEmailCount   uint8
	PasswordResetKey          string
	PasswordResetKeyExpiresAt time.Time
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}
