package dummy

import (
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

var Posting1 = model.Posting{
	ID:         1,
	UserName:   User1.Name,
	Title:      "This is a sample posting.",
	ImageURL:   "http://localhost:9000/toebeans-postings/20200101000000_testUser1",
	LikedCount: 0,
}

var Posting2 = model.Posting{
	ID:         2,
	UserName:   User2.Name,
	Title:      "test title",
	ImageURL:   "test url",
	LikedCount: 0,
}
