package model

import "time"

const MaxLimitPasswordResetPerDay = 3

type User struct {
	Name                      string
	Email                     string
	Password                  string
	Icon                      string
	SelfIntroduction          string
	ActivationKey             string
	EmailVerified             bool
	PasswordResetEmailCount   uint8
	PasswordResetKey          string
	PasswordResetKeyExpiresAt time.Time
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}
