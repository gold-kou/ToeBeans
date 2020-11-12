package usecase

import (
	"context"
	"os"

	"github.com/gold-kou/ToeBeans/app/adapter/aws"
	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

type DeletePostingUseCaseInterface interface {
	DeletePostingUseCase() error
}

type DeletePosting struct {
	ctx         context.Context
	tx          mysql.DBTransaction
	postingID   int64
	userName    string
	postingRepo *repository.PostingRepository
}

func NewDeletePosting(ctx context.Context, tx mysql.DBTransaction, postingID int64, userName string, postingRepo *repository.PostingRepository) *DeletePosting {
	return &DeletePosting{
		ctx:         ctx,
		tx:          tx,
		postingID:   postingID,
		userName:    userName,
		postingRepo: postingRepo,
	}
}

func (posting *DeletePosting) DeletePostingUseCase() error {
	p, err := posting.postingRepo.GetWhereIDUserName(posting.ctx, posting.postingID, posting.userName)
	if err != nil {
		return err
	}

	err = aws.DeleteObject(os.Getenv("S3_BUCKET_POSTINGS"), p.Title)
	if err != nil {
		return err
	}

	err = posting.tx.Do(posting.ctx, func(ctx context.Context) error {
		err := posting.postingRepo.DeleteWhereID(ctx, posting.postingID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
