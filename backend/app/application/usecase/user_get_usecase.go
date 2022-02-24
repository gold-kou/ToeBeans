package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

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
}

func NewGetUser(ctx context.Context, tx mysql.DBTransaction, tokenUserName, userName string, userRepo *repository.UserRepository) *GetUser {
	return &GetUser{
		ctx:           ctx,
		tx:            tx,
		tokenUserName: tokenUserName,
		userName:      userName,
		userRepo:      userRepo,
	}
}

func (user *GetUser) GetUserUseCase() (u model.User, err error) {
	// check userName in token exists
	_, err = user.userRepo.GetUserWhereName(user.ctx, user.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return model.User{}, lib.ErrTokenInvalidNotExistingUserName
		}
		return model.User{}, err
	}

	u, err = user.userRepo.GetUserWhereName(user.ctx, user.userName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return model.User{}, ErrNotExistsData
		}
		return
	}
	return
}
