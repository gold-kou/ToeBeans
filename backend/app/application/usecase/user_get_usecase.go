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
	tx             mysql.DBTransaction
	tokenUserName  string
	targetUserName string
	userRepo       *repository.UserRepository
	postingRepo    *repository.PostingRepository
	likeRepo       *repository.LikeRepository
	followRepo     *repository.FollowRepository
}

func NewGetUser(tx mysql.DBTransaction, tokenUserName, targetUserName string, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository, likeRepo *repository.LikeRepository, followRepo *repository.FollowRepository) *GetUser {
	return &GetUser{
		tx:             tx,
		tokenUserName:  tokenUserName,
		targetUserName: targetUserName,
		userRepo:       userRepo,
		postingRepo:    postingRepo,
		likeRepo:       likeRepo,
		followRepo:     followRepo,
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

	u, err = user.userRepo.GetUserWhereName(ctx, user.targetUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			err = ErrNotExistsData
		}
		return
	}

	postingCount, err = user.postingRepo.GetCountWhereUserID(ctx, u.ID)
	if err != nil {
		return
	}

	likeCount, err = user.likeRepo.GetLikeCountWhereUserID(ctx, u.ID)
	if err != nil {
		return
	}

	likedCount, err = user.likeRepo.GetLikedCountWhereUserID(ctx, u.ID)
	if err != nil {
		return
	}

	followCount, err = user.followRepo.GetFollowCountWhereUserID(ctx, u.ID)
	if err != nil {
		return
	}

	followedCount, err = user.followRepo.GetFollowedCountWhereUserID(ctx, u.ID)
	if err != nil {
		return
	}
	return
}
