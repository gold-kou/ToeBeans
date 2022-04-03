package usecase

import (
	"context"
	"errors"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

var ErrDeleteNotExistsFollow = errors.New("can't delete not existing follow")

type DeleteFollowUseCaseInterface interface {
	DeleteFollowUseCase() (*model.Follow, error)
}

type DeleteFollow struct {
	tx               mysql.DBTransaction
	followUserName   string
	followedUserName string
	userRepo         *repository.UserRepository
	followRepo       *repository.FollowRepository
}

func NewDeleteFollow(tx mysql.DBTransaction, followUserName, followedUserName string, userRepo *repository.UserRepository, followRepo *repository.FollowRepository) *DeleteFollow {
	return &DeleteFollow{
		tx:               tx,
		followUserName:   followUserName,
		followedUserName: followedUserName,
		userRepo:         userRepo,
		followRepo:       followRepo,
	}
}

func (follow *DeleteFollow) DeleteFollowUseCase(ctx context.Context) error {
	// check userName in token exists
	followingUser, err := follow.userRepo.GetUserWhereName(ctx, follow.followUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	// 存在しないユーザのフォロー削除はConflictエラー
	followedUser, err := follow.userRepo.GetUserWhereName(ctx, follow.followedUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrNotExitsUser
		}
		return err
	}

	// 存在しないフォローの削除はConflictエラー
	_, err = follow.followRepo.FindByBothUserIDs(ctx, followingUser.ID, followedUser.ID)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrDeleteNotExistsFollow
		}
		return err
	}

	err = follow.tx.Do(ctx, func(ctx context.Context) error {
		err := follow.followRepo.DeleteWhereBothUserIDs(ctx, followingUser.ID, followedUser.ID)
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
