package dummy

import (
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

var Follow1 = model.Follow{
	ID:                1,
	FollowingUserName: User1.Name,
	FollowedUserName:  User2.Name,
}
