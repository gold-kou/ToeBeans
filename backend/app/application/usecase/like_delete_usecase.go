package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type DeleteLikeUseCaseInterface interface {
	DeleteLikeUseCase() (*model.Like, error)
}

type DeleteLike struct {
	ctx           context.Context
	tx            mysql.DBTransaction
	tokenUserName string
	likeID        int64
	userRepo      *repository.UserRepository
	postingRepo   *repository.PostingRepository
	likeRepo      *repository.LikeRepository
}

func NewDeleteLike(ctx context.Context, tx mysql.DBTransaction, tokenUserName string, likeID int64, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository, likeRepo *repository.LikeRepository) *DeleteLike {
	return &DeleteLike{
		ctx:           ctx,
		tx:            tx,
		tokenUserName: tokenUserName,
		likeID:        likeID,
		userRepo:      userRepo,
		postingRepo:   postingRepo,
		likeRepo:      likeRepo,
	}
}

func (like *DeleteLike) DeleteLikeUseCase() error {
	// check userName in token exists
	_, err := like.userRepo.GetUserWhereName(like.ctx, like.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return lib.ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	l, err := like.likeRepo.GetWhereID(like.ctx, like.likeID)
	if err != nil {
		return err
	}

	err = like.tx.Do(like.ctx, func(ctx context.Context) error {
		if err := like.likeRepo.DeleteWhereID(ctx, like.likeID); err != nil {
			return err
		}

		// decrement
		if err := like.userRepo.UpdateLikeCount(ctx, like.tokenUserName, false); err != nil {
			return err
		}
		if err := like.userRepo.UpdateLikedCount(ctx, l.PostingID, false); err != nil {
			return err
		}
		if err := like.postingRepo.UpdateLikedCount(ctx, l.PostingID, false); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
