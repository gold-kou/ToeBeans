package dummy

import (
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

var Posting1 = model.Posting{
	ID:       1,
	UserID:   User1.ID,
	Title:    "This is a sample posting.",
	ImageURL: "http://localhost:9000/toebeans-postings/20200101000000_testUser1",
}

var Posting2 = model.Posting{
	ID:       2,
	UserID:   User2.ID,
	Title:    "test title",
	ImageURL: "test url",
}
