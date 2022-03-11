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
	ctx              context.Context
	tx               mysql.DBTransaction
	followUserName   string
	followedUserName string
	userRepo         *repository.UserRepository
	followRepo       *repository.FollowRepository
}

func NewDeleteFollow(ctx context.Context, tx mysql.DBTransaction, followUserName, followedUserName string, userRepo *repository.UserRepository, followRepo *repository.FollowRepository) *DeleteFollow {
	return &DeleteFollow{
		ctx:              ctx,
		tx:               tx,
		followUserName:   followUserName,
		followedUserName: followedUserName,
		userRepo:         userRepo,
		followRepo:       followRepo,
	}
}

func (follow *DeleteFollow) DeleteFollowUseCase() error {
	// check userName in token exists
	_, err := follow.userRepo.GetUserWhereName(follow.ctx, follow.followUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	// 存在しないユーザのフォロー削除はConflictエラー
	_, err = follow.userRepo.GetUserWhereName(follow.ctx, follow.followedUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrNotExitsUser
		}
		return err
	}

	// 存在しないフォローの削除はConflictエラー
	_, err = follow.followRepo.FindByBothUserNames(follow.ctx, follow.followUserName, follow.followedUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrDeleteNotExistsFollow
		}
		return err
	}

	err = follow.tx.Do(follow.ctx, func(ctx context.Context) error {
		err := follow.followRepo.DeleteWhereFollowingFollowedUserName(ctx, follow.followUserName, follow.followedUserName)
		if err != nil {
			return err
		}

		err = follow.userRepo.UpdateFollowCount(ctx, follow.followUserName, false)
		if err != nil {
			return err
		}

		err = follow.userRepo.UpdateFollowedCount(ctx, follow.followedUserName, false)
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
