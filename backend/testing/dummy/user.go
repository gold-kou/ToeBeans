package dummy

import (
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

var User1 = model.User{
	Name:             "testUser1",
	Email:            "testUser1@example.com",
	Password:         "Password1234", // raw: password
	ActivationKey:    "b2add849-2c94-4a4f-a079-baf447536bc0",
	PasswordResetKey: "fec668ed-f69e-45cd-87b9-a76f27759134",
}

var User2 = model.User{
	Name:     "testUser2",
	Email:    "testUser2@example.com",
	Password: "Password5678", // raw: password
}

var SecretKey = "test_secret_key"
