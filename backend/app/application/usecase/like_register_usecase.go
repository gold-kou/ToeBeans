package usecase

import (
	"context"
	"errors"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

var ErrLikeYourPosting = errors.New("you can't like your posting")
var ErrAlreadyLiked = errors.New("Whoops, you already liked the posting")

type RegisterLikeUseCaseInterface interface {
	RegisterLikeUseCase() (*model.Like, error)
}

type RegisterLike struct {
	tx               mysql.DBTransaction
	tokenUserID      int64
	tokenUserName    string
	postingID        int
	userRepo         *repository.UserRepository
	postingRepo      *repository.PostingRepository
	likeRepo         *repository.LikeRepository
	notificationRepo *repository.NotificationRepository
}

func NewRegisterLike(tx mysql.DBTransaction, tokenUserID int64, tokenUserName string, postingID int, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository, likeRepo *repository.LikeRepository, notificationRepo *repository.NotificationRepository) *RegisterLike {
	return &RegisterLike{
		tx:               tx,
		tokenUserID:      tokenUserID,
		tokenUserName:    tokenUserName,
		postingID:        postingID,
		userRepo:         userRepo,
		postingRepo:      postingRepo,
		likeRepo:         likeRepo,
		notificationRepo: notificationRepo,
	}
}

func (like *RegisterLike) RegisterLikeUseCase(ctx context.Context) error {
	// check userName in token exists
	_, err := like.userRepo.GetUserWhereName(ctx, like.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	p, err := like.postingRepo.GetWhereID(ctx, int64(like.postingID))
	if err != nil {
		return err
	}

	if like.tokenUserID == p.UserID {
		return ErrLikeYourPosting
	}

	err = like.tx.Do(ctx, func(ctx context.Context) error {
		l := model.Like{
			UserID:    like.tokenUserID,
			PostingID: int64(like.postingID),
		}
		if err := like.likeRepo.Create(ctx, &l); err != nil {
			if err == repository.ErrDuplicateData {
				return ErrAlreadyLiked
			}
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
