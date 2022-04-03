package model

import "time"

const MaxLimitPasswordResetPerDay = 3

type PasswordReset struct {
	ID                        int64
	UserID                    int64
	PasswordResetEmailCount   uint8
	PasswordResetKey          string
	PasswordResetKeyExpiresAt time.Time
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}
