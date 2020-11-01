package dummy

import "github.com/gold-kou/ToeBeans/app/domain/model"

var User1 = model.User{
	Name:     "test1",
	Email:    "test1@example.com",
	Password: "Password1234", // raw: password
}

var SecretKey = "test_secret_key"
