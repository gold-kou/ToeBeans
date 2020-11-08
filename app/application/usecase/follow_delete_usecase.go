package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/model"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

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
	err := follow.tx.Do(follow.ctx, func(ctx context.Context) error {
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
