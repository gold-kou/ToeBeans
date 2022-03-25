package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/aws"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type DeletePostingUseCaseInterface interface {
	DeletePostingUseCase() error
}

type DeletePosting struct {
	tx            mysql.DBTransaction
	postingID     int64
	tokenUserName string
	userRepo      *repository.UserRepository
	postingRepo   *repository.PostingRepository
}

func NewDeletePosting(tx mysql.DBTransaction, postingID int64, tokenUserName string, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository) *DeletePosting {
	return &DeletePosting{
		tx:            tx,
		postingID:     postingID,
		tokenUserName: tokenUserName,
		userRepo:      userRepo,
		postingRepo:   postingRepo,
	}
}

func (posting *DeletePosting) DeletePostingUseCase(ctx context.Context) error {
	// check userName in token exists
	_, err := posting.userRepo.GetUserWhereName(ctx, posting.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	p, err := posting.postingRepo.GetWhereIDUserName(ctx, posting.postingID, posting.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrNotExistsData
		}
		return err
	}

	err = aws.DeleteObject(bucketPosting, p.Title)
	if err != nil {
		return err
	}

	err = posting.tx.Do(ctx, func(ctx context.Context) error {
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
