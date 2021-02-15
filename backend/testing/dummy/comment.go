package dummy

import (
	"github.com/gold-kou/ToeBeans/app/domain/model"
)

var Comment1 = model.Comment{
	ID:        1,
	UserName:  User1.Name,
	PostingID: Posting1.ID,
	Comment:   "test comment",
}
