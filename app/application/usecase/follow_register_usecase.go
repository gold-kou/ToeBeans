package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/model"
	modelHTTP "github.com/gold-kou/ToeBeans/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

type RegisterFollowUseCaseInterface interface {
	RegisterFollowUseCase() (*model.Follow, error)
}

type RegisterFollow struct {
	ctx               context.Context
	tx                mysql.DBTransaction
	userName          string
	reqRegisterFollow *modelHTTP.Follow
	userRepo          *repository.UserRepository
	followRepo        *repository.FollowRepository
	notificationRepo  *repository.NotificationRepository
}

func NewRegisterFollow(ctx context.Context, tx mysql.DBTransaction, userName string, reqRegisterFollow *modelHTTP.Follow, userRepo *repository.UserRepository, followRepo *repository.FollowRepository, notificationRepo *repository.NotificationRepository) *RegisterFollow {
	return &RegisterFollow{
		ctx:               ctx,
		tx:                tx,
		userName:          userName,
		reqRegisterFollow: reqRegisterFollow,
		userRepo:          userRepo,
		followRepo:        followRepo,
		notificationRepo:  notificationRepo,
	}
}

func (follow *RegisterFollow) RegisterFollowUseCase() error {
	err := follow.tx.Do(follow.ctx, func(ctx context.Context) error {
		u := model.Follow{
			FollowingUserName: follow.userName,
			FollowedUserName:  follow.reqRegisterFollow.FollowedUserName,
		}
		err := follow.followRepo.Create(ctx, &u)
		if err != nil {
			return err
		}

		err = follow.userRepo.UpdateFollowCount(ctx, follow.userName, true)
		if err != nil {
			return err
		}

		err = follow.userRepo.UpdateFollowedCount(ctx, follow.reqRegisterFollow.FollowedUserName, true)
		if err != nil {
			return err
		}

		// TODO notification
		// if follow.userName != follow.reqRegisterFollow.FollowedUserName {
		// 	n := model.Notification{
		// 		VisitorName: follow.userName,
		// 		VisitedName: follow.reqRegisterFollow.FollowedUserName,
		// 		Action:      model.FollowAction,
		// 	}
		// 	if err = follow.notificationRepo.Create(ctx, &n); err != nil {
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
