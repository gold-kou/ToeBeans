package usecase

import (
	"context"
	"errors"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

var ErrDeleteNotExistsLike = errors.New("can't delete not existing like")

type DeleteLikeUseCaseInterface interface {
	DeleteLikeUseCase() (*model.Like, error)
}

type DeleteLike struct {
	ctx           context.Context
	tx            mysql.DBTransaction
	tokenUserName string
	postingID     int64
	userRepo      *repository.UserRepository
	postingRepo   *repository.PostingRepository
	likeRepo      *repository.LikeRepository
}

func NewDeleteLike(ctx context.Context, tx mysql.DBTransaction, tokenUserName string, postingID int64, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository, likeRepo *repository.LikeRepository) *DeleteLike {
	return &DeleteLike{
		ctx:           ctx,
		tx:            tx,
		tokenUserName: tokenUserName,
		postingID:     postingID,
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
			return ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	_, err = like.likeRepo.GetWhereUserNamePostingID(like.ctx, like.tokenUserName, like.postingID)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrDeleteNotExistsLike
		}
		return err
	}

	err = like.tx.Do(like.ctx, func(ctx context.Context) error {
		if err := like.likeRepo.DeleteWhereUserNamePostingID(ctx, like.tokenUserName, like.postingID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
