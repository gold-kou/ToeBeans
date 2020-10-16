package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/model"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

type GetUserUseCaseInterface interface {
	GetUserUseCase() (*model.User, error)
}

type GetUser struct {
	ctx      context.Context
	tx       mysql.DBTransaction
	userName string
	userRepo *repository.UserRepository
}

func NewGetUser(ctx context.Context, tx mysql.DBTransaction, userName string, userRepo *repository.UserRepository) *GetUser {
	return &GetUser{
		ctx:      ctx,
		tx:       tx,
		userName: userName,
		userRepo: userRepo,
	}
}

func (user *GetUser) GetUserUseCase() (u model.User, err error) {
	u, err = user.userRepo.GetUserWhereName(user.ctx, user.userName)
	if err != nil {
		return
	}
	return
}
