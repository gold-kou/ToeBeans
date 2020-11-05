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
	ctx      context.Context
	tx       mysql.DBTransaction
	userName string
	userRepo *repository.UserRepository
}

func NewDeleteUser(ctx context.Context, tx mysql.DBTransaction, userName string, userRepo *repository.UserRepository) *DeleteUser {
	return &DeleteUser{
		ctx:      ctx,
		tx:       tx,
		userName: userName,
		userRepo: userRepo,
	}
}

func (user *DeleteUser) DeleteUserUseCase() error {
	err := user.tx.Do(user.ctx, func(ctx context.Context) error {
		// TODO notification→likes/comments/follows→postings→usersで削除
		err := user.userRepo.DeleteWhereName(ctx, user.userName)
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
