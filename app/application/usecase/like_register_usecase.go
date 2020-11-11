package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/model"
	modelHTTP "github.com/gold-kou/ToeBeans/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

type RegisterLikeUseCaseInterface interface {
	RegisterLikeUseCase() (*model.Like, error)
}

type RegisterLike struct {
	ctx              context.Context
	tx               mysql.DBTransaction
	userName         string
	reqRegisterLike  *modelHTTP.Like
	userRepo         *repository.UserRepository
	postingRepo      *repository.PostingRepository
	likeRepo         *repository.LikeRepository
	notificationRepo *repository.NotificationRepository
}

func NewRegisterLike(ctx context.Context, tx mysql.DBTransaction, userName string, reqRegisterLike *modelHTTP.Like, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository, likeRepo *repository.LikeRepository, notificationRepo *repository.NotificationRepository) *RegisterLike {
	return &RegisterLike{
		ctx:              ctx,
		tx:               tx,
		userName:         userName,
		reqRegisterLike:  reqRegisterLike,
		userRepo:         userRepo,
		postingRepo:      postingRepo,
		likeRepo:         likeRepo,
		notificationRepo: notificationRepo,
	}
}

func (like *RegisterLike) RegisterLikeUseCase() error {
	p, err := like.postingRepo.GetWhereID(like.ctx, like.reqRegisterLike.PostingId)
	if err != nil {
		return err
	}

	if like.userName == p.UserName {
		return ErrLikeYourSelf
	}

	err = like.tx.Do(like.ctx, func(ctx context.Context) error {
		l := model.Like{
			UserName:  like.userName,
			PostingID: like.reqRegisterLike.PostingId,
		}
		if err := like.likeRepo.Create(ctx, &l); err != nil {
			return err
		}

		// increment
		if err := like.userRepo.UpdateLikeCount(ctx, like.userName, true); err != nil {
			return err
		}
		if err := like.userRepo.UpdateLikedCount(ctx, like.reqRegisterLike.PostingId, true); err != nil {
			return err
		}
		if err := like.postingRepo.UpdateLikedCount(ctx, like.reqRegisterLike.PostingId, true); err != nil {
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
