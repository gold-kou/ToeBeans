package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/model"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

type DeleteLikeUseCaseInterface interface {
	DeleteLikeUseCase() (*model.Like, error)
}

type DeleteLike struct {
	ctx         context.Context
	tx          mysql.DBTransaction
	userName    string
	likeID      int64
	userRepo    *repository.UserRepository
	postingRepo *repository.PostingRepository
	likeRepo    *repository.LikeRepository
}

func NewDeleteLike(ctx context.Context, tx mysql.DBTransaction, userName string, likeID int64, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository, likeRepo *repository.LikeRepository) *DeleteLike {
	return &DeleteLike{
		ctx:         ctx,
		tx:          tx,
		userName:    userName,
		likeID:      likeID,
		userRepo:    userRepo,
		postingRepo: postingRepo,
		likeRepo:    likeRepo,
	}
}

func (like *DeleteLike) DeleteLikeUseCase() error {
	l, err := like.likeRepo.GetWhereID(like.ctx, like.likeID)
	if err != nil {
		return err
	}

	err = like.tx.Do(like.ctx, func(ctx context.Context) error {
		if err := like.likeRepo.DeleteWhereID(ctx, like.likeID); err != nil {
			return err
		}

		// decrement
		if err := like.userRepo.UpdateLikeCount(ctx, like.userName, false); err != nil {
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
