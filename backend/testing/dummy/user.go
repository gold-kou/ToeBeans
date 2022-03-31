package dummy

import (
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

var User1 = model.User{
	ID:            1,
	Name:          "testUser1",
	Email:         "testUser1@example.com",
	Password:      "Password1234",
	ActivationKey: "b2add849-2c94-4a4f-a079-baf447536bc0",
}

var User2 = model.User{
	ID:       2,
	Name:     "testUser2",
	Email:    "testUser2@example.com",
	Password: "Password5678",
}

var User3 = model.User{
	ID:       3,
	Name:     "testUser3",
	Email:    "testUser3@example.com",
	Password: "Password3333",
}

var SecretKey = "test_secret_key"
