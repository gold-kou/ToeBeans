package dummy

import (
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

var Like1to2 = model.Like{
	ID:        1,
	UserID:    User1.ID,
	PostingID: Posting2.ID, // you can't like yourself posting
}

var Like2to1 = model.Like{
	ID:        2,
	UserID:    User2.ID,
	PostingID: Posting1.ID, // you can't like yourself posting
}
