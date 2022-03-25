package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type GetUserUseCaseInterface interface {
	GetUserUseCase() (*model.User, error)
}

type GetUser struct {
	tx            mysql.DBTransaction
	tokenUserName string
	userName      string
	userRepo      *repository.UserRepository
	postingRepo   *repository.PostingRepository
	likeRepo      *repository.LikeRepository
	followRepo    *repository.FollowRepository
}

func NewGetUser(tx mysql.DBTransaction, tokenUserName, userName string, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository, likeRepo *repository.LikeRepository, followRepo *repository.FollowRepository) *GetUser {
	return &GetUser{
		tx:            tx,
		tokenUserName: tokenUserName,
		userName:      userName,
		userRepo:      userRepo,
		postingRepo:   postingRepo,
		likeRepo:      likeRepo,
		followRepo:    followRepo,
	}
}

func (user *GetUser) GetUserUseCase(ctx context.Context) (u model.User, postingCount, likeCount, likedCount, followCount, followedCount int64, err error) {
	// check userName in token exists
	_, err = user.userRepo.GetUserWhereName(ctx, user.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			err = ErrTokenInvalidNotExistingUserName
		}
		return
	}

	u, err = user.userRepo.GetUserWhereName(ctx, user.userName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			err = ErrNotExistsData
		}
		return
	}

	postingCount, err = user.postingRepo.GetCountWhereUserName(ctx, user.userName)
	if err != nil {
		return
	}

	likeCount, err = user.likeRepo.GetLikeCountWhereUserName(ctx, user.userName)
	if err != nil {
		return
	}

	likedCount, err = user.likeRepo.GetLikedCountWhereUserName(ctx, user.userName)
	if err != nil {
		return
	}

	followCount, err = user.followRepo.GetFollowCountWhereUserName(ctx, user.userName)
	if err != nil {
		return
	}

	followedCount, err = user.followRepo.GetFollowedCountWhereUserName(ctx, user.userName)
	if err != nil {
		return
	}
	return
}
