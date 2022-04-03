package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/domain/model"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type DeleteUserUseCaseInterface interface {
	DeleteUserUseCase() (*model.User, error)
}

type DeleteUser struct {
	tx                mysql.DBTransaction
	userName          string
	userRepo          *repository.UserRepository
	passwordResetRepo *repository.PasswordResetRepository
	postingRepo       *repository.PostingRepository
	likeRepo          *repository.LikeRepository
	commentRepo       *repository.CommentRepository
	followRepo        *repository.FollowRepository
}

func NewDeleteUser(tx mysql.DBTransaction, userName string, userRepo *repository.UserRepository, passwordResetRepo *repository.PasswordResetRepository, postingRepo *repository.PostingRepository, likeRepo *repository.LikeRepository, commentRepo *repository.CommentRepository, followRepo *repository.FollowRepository) *DeleteUser {
	return &DeleteUser{
		tx:                tx,
		userName:          userName,
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
		postingRepo:       postingRepo,
		likeRepo:          likeRepo,
		commentRepo:       commentRepo,
		followRepo:        followRepo,
	}
}

func (user *DeleteUser) DeleteUserUseCase(ctx context.Context) error {
	// check user exists
	u, err := user.userRepo.GetUserWhereName(ctx, user.userName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrNotExitsUser
		}
		return err
	}

	err = user.tx.Do(ctx, func(ctx context.Context) error {
		// TODO notification delete

		err = user.likeRepo.DeleteWhereUserID(ctx, u.ID)
		if err != nil {
			return err
		}

		// 削除対象ユーザの投稿に対するいいねを削除する
		err = user.likeRepo.DeleteWhereInPosingIDs(ctx, u.ID)
		if err != nil {
			return err
		}

		err = user.commentRepo.DeleteWhereUserID(ctx, u.ID)
		if err != nil {
			return err
		}

		err = user.followRepo.DeleteWhereFollowingUserID(ctx, u.ID)
		if err != nil {
			return err
		}

		err = user.followRepo.DeleteWhereFollowedUserID(ctx, u.ID)
		if err != nil {
			return err
		}

		err = user.postingRepo.DeleteWhereUserID(ctx, u.ID)
		if err != nil {
			return err
		}

		err = user.passwordResetRepo.DeleteWhereUserID(ctx, u.ID)
		if err != nil {
			return err
		}

		err = user.userRepo.DeleteWhereID(ctx, u.ID)
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
