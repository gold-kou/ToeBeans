package dummy

import (
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

var PostingReport1 = model.PostingReport{
	ID:        1,
	PostingID: int(Posting1.ID),
	Detail:    "an inappropriate posting",
}
