package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type RegisterLikeUseCaseInterface interface {
	RegisterLikeUseCase() (*model.Like, error)
}

type RegisterLike struct {
	ctx              context.Context
	tx               mysql.DBTransaction
	tokenUserName    string
	postingID        int
	userRepo         *repository.UserRepository
	postingRepo      *repository.PostingRepository
	likeRepo         *repository.LikeRepository
	notificationRepo *repository.NotificationRepository
}

func NewRegisterLike(ctx context.Context, tx mysql.DBTransaction, tokenUserName string, postingID int, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository, likeRepo *repository.LikeRepository, notificationRepo *repository.NotificationRepository) *RegisterLike {
	return &RegisterLike{
		ctx:              ctx,
		tx:               tx,
		tokenUserName:    tokenUserName,
		postingID:        postingID,
		userRepo:         userRepo,
		postingRepo:      postingRepo,
		likeRepo:         likeRepo,
		notificationRepo: notificationRepo,
	}
}

func (like *RegisterLike) RegisterLikeUseCase() error {
	// check userName in token exists
	_, err := like.userRepo.GetUserWhereName(like.ctx, like.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return lib.ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	p, err := like.postingRepo.GetWhereID(like.ctx, int64(like.postingID))
	if err != nil {
		return err
	}

	if like.tokenUserName == p.UserName {
		return ErrLikeYourSelf
	}

	err = like.tx.Do(like.ctx, func(ctx context.Context) error {
		l := model.Like{
			UserName:  like.tokenUserName,
			PostingID: int64(like.postingID),
		}
		if err := like.likeRepo.Create(ctx, &l); err != nil {
			return err
		}

		// increment
		if err := like.userRepo.UpdateLikeCount(ctx, like.tokenUserName, true); err != nil {
			return err
		}
		if err := like.userRepo.UpdateLikedCount(ctx, int64(like.postingID), true); err != nil {
			return err
		}
		if err := like.postingRepo.UpdateLikedCount(ctx, int64(like.postingID), true); err != nil {
			return err
		}

		// TODO notification
		// if like.userName != p.UserName {
		// 	n := model.Notification{
		// 		VisitorName: like.userName,
		// 		VisitedName: p.UserName,
		// 		Action:      model.LikeAction,
		// 	}
		// 	if err = like.notificationRepo.Create(ctx, &n); err != nil {
		// 		return err
		// 	}
		// }

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
