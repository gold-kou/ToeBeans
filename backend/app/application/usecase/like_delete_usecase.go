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
	tx            mysql.DBTransaction
	tokenUserName string
	postingID     int64
	userRepo      *repository.UserRepository
	postingRepo   *repository.PostingRepository
	likeRepo      *repository.LikeRepository
}

func NewDeleteLike(tx mysql.DBTransaction, tokenUserName string, postingID int64, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository, likeRepo *repository.LikeRepository) *DeleteLike {
	return &DeleteLike{
		tx:            tx,
		tokenUserName: tokenUserName,
		postingID:     postingID,
		userRepo:      userRepo,
		postingRepo:   postingRepo,
		likeRepo:      likeRepo,
	}
}

func (like *DeleteLike) DeleteLikeUseCase(ctx context.Context) error {
	// check userName in token exists
	_, err := like.userRepo.GetUserWhereName(ctx, like.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	_, err = like.likeRepo.GetWhereUserNamePostingID(ctx, like.tokenUserName, like.postingID)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrDeleteNotExistsLike
		}
		return err
	}

	err = like.tx.Do(ctx, func(ctx context.Context) error {
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
