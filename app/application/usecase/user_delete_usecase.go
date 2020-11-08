package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/app/domain/model"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

type DeleteUserUseCaseInterface interface {
	DeleteUserUseCase() (*model.User, error)
}

type DeleteUser struct {
	ctx         context.Context
	tx          mysql.DBTransaction
	userName    string
	userRepo    *repository.UserRepository
	postingRepo *repository.PostingRepository
	likeRepo    *repository.LikeRepository
	commentRepo *repository.CommentRepository
	followRepo  *repository.FollowRepository
}

func NewDeleteUser(ctx context.Context, tx mysql.DBTransaction, userName string, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository, likeRepo *repository.LikeRepository, commentRepo *repository.CommentRepository, followRepo *repository.FollowRepository) *DeleteUser {
	return &DeleteUser{
		ctx:         ctx,
		tx:          tx,
		userName:    userName,
		userRepo:    userRepo,
		postingRepo: postingRepo,
		likeRepo:    likeRepo,
		commentRepo: commentRepo,
		followRepo:  followRepo,
	}
}

func (user *DeleteUser) DeleteUserUseCase() error {
	err := user.tx.Do(user.ctx, func(ctx context.Context) error {
		// TODO notification削除
		err := user.likeRepo.DeleteWhereUserName(ctx, user.userName)
		if err != nil {
			return err
		}

		err = user.commentRepo.DeleteWhereUserName(ctx, user.userName)
		if err != nil {
			return err
		}

		err = user.followRepo.DeleteWhereFollowingUserName(ctx, user.userName)
		if err != nil {
			return err
		}

		err = user.postingRepo.DeleteWhereUserName(ctx, user.userName)
		if err != nil {
			return err
		}

		err = user.userRepo.DeleteWhereName(ctx, user.userName)
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
