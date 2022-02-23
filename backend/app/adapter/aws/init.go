package aws

import (
	"os"

	"github.com/gold-kou/ToeBeans/backend/app"
)

var accessKey, secretKey, region string
var emailFrom string

func init() {
	accessKey = os.Getenv("AWS_ACCESS_KEY")
	secretKey = os.Getenv("AWS_SECRET_KEY")
	region = os.Getenv("AWS_REGION")
	if (accessKey == "" || secretKey == "" || region == "") && (app.IsLocal() || app.IsTest()) {
		panic("something aws environment is unset")
	}

	emailFrom = os.Getenv("SYSTEM_EMAIL")
	if emailFrom == "" {
		panic("SYSTEM_EMAIL is unset")
	}
}
