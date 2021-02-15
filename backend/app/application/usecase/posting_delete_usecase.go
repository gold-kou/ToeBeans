package usecase

import (
	"context"
	"os"

	"github.com/gold-kou/ToeBeans/app/lib"

	"github.com/gold-kou/ToeBeans/app/adapter/aws"
	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

type DeletePostingUseCaseInterface interface {
	DeletePostingUseCase() error
}

type DeletePosting struct {
	ctx           context.Context
	tx            mysql.DBTransaction
	postingID     int64
	tokenUserName string
	userRepo      *repository.UserRepository
	postingRepo   *repository.PostingRepository
}

func NewDeletePosting(ctx context.Context, tx mysql.DBTransaction, postingID int64, tokenUserName string, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository) *DeletePosting {
	return &DeletePosting{
		ctx:           ctx,
		tx:            tx,
		postingID:     postingID,
		tokenUserName: tokenUserName,
		userRepo:      userRepo,
		postingRepo:   postingRepo,
	}
}

func (posting *DeletePosting) DeletePostingUseCase() error {
	// check userName in token exists
	_, err := posting.userRepo.GetUserWhereName(posting.ctx, posting.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return lib.ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	p, err := posting.postingRepo.GetWhereIDUserName(posting.ctx, posting.postingID, posting.tokenUserName)
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
