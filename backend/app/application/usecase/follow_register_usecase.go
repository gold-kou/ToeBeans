package usecase

import (
	"context"

	"github.com/pkg/errors"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

var ErrFollowedUserNotExists = errors.New("the followed user doesn't exsist")
var ErrAlreadyFollowed = errors.New("the user is already followed by you")

type RegisterFollowUseCaseInterface interface {
	RegisterFollowUseCase() (*model.Follow, error)
}

type RegisterFollow struct {
	tx               mysql.DBTransaction
	tokenUserID      int64
	tokenUserName    string
	followedUserName string
	userRepo         *repository.UserRepository
	followRepo       *repository.FollowRepository
	notificationRepo *repository.NotificationRepository
}

func NewRegisterFollow(tx mysql.DBTransaction, tokenUserID int64, tokenUserName string, followedUserName string, userRepo *repository.UserRepository, followRepo *repository.FollowRepository, notificationRepo *repository.NotificationRepository) *RegisterFollow {
	return &RegisterFollow{
		tx:               tx,
		tokenUserID:      tokenUserID,
		tokenUserName:    tokenUserName,
		followedUserName: followedUserName,
		userRepo:         userRepo,
		followRepo:       followRepo,
		notificationRepo: notificationRepo,
	}
}

func (follow *RegisterFollow) RegisterFollowUseCase(ctx context.Context) error {
	// check userName in token exists
	_, err := follow.userRepo.GetUserWhereName(ctx, follow.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	followedUser, err := follow.userRepo.GetUserWhereName(ctx, follow.followedUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrFollowedUserNotExists
		}
		return err
	}

	err = follow.tx.Do(ctx, func(ctx context.Context) error {
		u := model.Follow{
			FollowingUserID: follow.tokenUserID,
			FollowedUserID:  followedUser.ID,
		}
		err := follow.followRepo.Create(ctx, &u)
		if err != nil {
			if err == repository.ErrDuplicateData {
				return ErrAlreadyFollowed
			}
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
