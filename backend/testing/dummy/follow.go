package dummy

import (
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

var Follow1to2 = model.Follow{
	ID:                1,
	FollowingUserName: User1.Name,
	FollowedUserName:  User2.Name,
}

var Follow2to1 = model.Follow{
	ID:                2,
	FollowingUserName: User2.Name,
	FollowedUserName:  User1.Name,
}
