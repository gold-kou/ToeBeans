package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type GetFollowStateUseCaseInterface interface {
	GetFollowStateUseCase() (*model.Follow, error)
}

type GetFollowState struct {
	ctx              context.Context
	tx               mysql.DBTransaction
	tokenUserName    string
	followedUserName string
	userRepo         *repository.UserRepository
	followRepo       *repository.FollowRepository
}

func NewGetFollowState(ctx context.Context, tx mysql.DBTransaction, tokenUserName string, followedUserName string, userRepo *repository.UserRepository, followRepo *repository.FollowRepository) *GetFollowState {
	return &GetFollowState{
		ctx:              ctx,
		tx:               tx,
		tokenUserName:    tokenUserName,
		followedUserName: followedUserName,
		userRepo:         userRepo,
		followRepo:       followRepo,
	}
}

func (follow *GetFollowState) GetFollowStateUseCase() error {
	// check userName in token exists
	_, err := follow.userRepo.GetUserWhereName(follow.ctx, follow.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	_, err = follow.followRepo.FindByBothUserNames(follow.ctx, follow.tokenUserName, follow.followedUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrNotExistsData
		}
		return err
	}

	return nil
}
