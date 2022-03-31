package model

import "time"

type User struct {
	ID               int64
	Name             string
	Email            string
	Password         string
	Icon             string
	SelfIntroduction string
	ActivationKey    string
	EmailVerified    bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
