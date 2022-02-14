package app

import (
	"os"
	"strings"
)

var AppEnv string

func init() {
	AppEnv = strings.ToLower(os.Getenv("APP_ENV"))
	if AppEnv == "" {
		panic("APP_ENV is unset")
	}
}

func IsProduction() bool {
	return AppEnv == "prd"
}

func IsLocal() bool {
	return AppEnv == "local"
}

func IsTest() bool {
	return AppEnv == "test"
}
