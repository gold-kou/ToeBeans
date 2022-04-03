package dummy

import (
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

var Comment1 = model.Comment{
	ID:        1,
	UserID:    User1.ID,
	PostingID: Posting1.ID,
	Comment:   "test comment",
}
