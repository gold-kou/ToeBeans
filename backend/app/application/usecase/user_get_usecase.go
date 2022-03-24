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
	ctx           context.Context
	tx            mysql.DBTransaction
	tokenUserName string
	userName      string
	userRepo      *repository.UserRepository
	postingRepo   *repository.PostingRepository
	likeRepo      *repository.LikeRepository
	followRepo    *repository.FollowRepository
}

func NewGetUser(ctx context.Context, tx mysql.DBTransaction, tokenUserName, userName string, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository, likeRepo *repository.LikeRepository, followRepo *repository.FollowRepository) *GetUser {
	return &GetUser{
		ctx:           ctx,
		tx:            tx,
		tokenUserName: tokenUserName,
		userName:      userName,
		userRepo:      userRepo,
		postingRepo:   postingRepo,
		likeRepo:      likeRepo,
		followRepo:    followRepo,
	}
}

func (user *GetUser) GetUserUseCase() (u model.User, postingCount, likeCount, likedCount, followCount, followedCount int64, err error) {
	// check userName in token exists
	_, err = user.userRepo.GetUserWhereName(user.ctx, user.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			err = ErrTokenInvalidNotExistingUserName
		}
		return
	}

	u, err = user.userRepo.GetUserWhereName(user.ctx, user.userName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			err = ErrNotExistsData
		}
		return
	}

	postingCount, err = user.postingRepo.GetCountWhereUserName(user.ctx, user.userName)
	if err != nil {
		return
	}

	likeCount, err = user.likeRepo.GetLikeCountWhereUserName(user.ctx, user.userName)
	if err != nil {
		return
	}

	likedCount, err = user.likeRepo.GetLikedCountWhereUserName(user.ctx, user.userName)
	if err != nil {
		return
	}

	followCount, err = user.followRepo.GetFollowCountWhereUserName(user.ctx, user.userName)
	if err != nil {
		return
	}

	followedCount, err = user.followRepo.GetFollowedCountWhereUserName(user.ctx, user.userName)
	if err != nil {
		return
	}
	return
}
