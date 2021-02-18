package dummy

import (
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

var Like1 = model.Like{
	ID:        1,
	UserName:  User1.Name,
	PostingID: Posting2.ID, // you can't like yourself posting
}
