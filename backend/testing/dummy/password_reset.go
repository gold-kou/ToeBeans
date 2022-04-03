package dummy

import "github.com/gold-kou/ToeBeans/backend/app/domain/model"

var PasswordReset1 = model.PasswordReset{
	ID:                      1,
	UserID:                  User1.ID,
	PasswordResetEmailCount: 1,
	PasswordResetKey:        "fec668ed-f69e-45cd-87b9-a76f27759134",
}
