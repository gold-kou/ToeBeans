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
	tx               mysql.DBTransaction
	tokenUserName    string
	followedUserName string
	userRepo         *repository.UserRepository
	followRepo       *repository.FollowRepository
}

func NewGetFollowState(tx mysql.DBTransaction, tokenUserName string, followedUserName string, userRepo *repository.UserRepository, followRepo *repository.FollowRepository) *GetFollowState {
	return &GetFollowState{
		tx:               tx,
		tokenUserName:    tokenUserName,
		followedUserName: followedUserName,
		userRepo:         userRepo,
		followRepo:       followRepo,
	}
}

func (follow *GetFollowState) GetFollowStateUseCase(ctx context.Context) error {
	// check userName in token exists
	_, err := follow.userRepo.GetUserWhereName(ctx, follow.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	_, err = follow.followRepo.FindByBothUserNames(ctx, follow.tokenUserName, follow.followedUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrNotExistsData
		}
		return err
	}

	return nil
}
