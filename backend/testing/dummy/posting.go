package dummy

import (
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

var Posting1 = model.Posting{
	ID:         1,
	UserName:   User1.Name,
	Title:      "20200101000000_testUser1_This is a sample posting.",
	ImageURL:   "http://minio:9000/postings/20200101000000_testUser1_This%20is%20a%20sample%20posting.",
	LikedCount: 0,
}

var Posting2 = model.Posting{
	ID:         2,
	UserName:   User2.Name,
	Title:      "test title",
	ImageURL:   "test url",
	LikedCount: 0,
}
