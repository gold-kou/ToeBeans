package dummy

import (
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

var Follow1to2 = model.Follow{
	ID:              1,
	FollowingUserID: User1.ID,
	FollowedUserID:  User2.ID,
}

var Follow2to1 = model.Follow{
	ID:              2,
	FollowingUserID: User2.ID,
	FollowedUserID:  User1.ID,
}
