package dummy

import (
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

var Like1to2 = model.Like{
	ID:        1,
	UserName:  User1.Name,
	PostingID: Posting2.ID, // you can't like yourself posting
}

var Like2to1 = model.Like{
	ID:        2,
	UserName:  User2.Name,
	PostingID: Posting1.ID, // you can't like yourself posting
}
