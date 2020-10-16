package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/app/domain/model"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

type UpdateActivationUseCaseInterface interface {
	UserActivationUseCase() (*model.User, error)
}

type UserActivation struct {
	ctx           context.Context
	tx            mysql.DBTransaction
	userName      string
	activationKey string
	userRepo      *repository.UserRepository
}

func NewUserActivation(ctx context.Context, tx mysql.DBTransaction, userName, activationKey string, userRepo *repository.UserRepository) *UserActivation {
	return &UserActivation{
		ctx:           ctx,
		tx:            tx,
		userName:      userName,
		activationKey: activationKey,
		userRepo:      userRepo,
	}
}

func (ua *UserActivation) UserActivationUseCase() error {
	err := ua.tx.Do(ua.ctx, func(ctx context.Context) error {
		err := ua.userRepo.UpdateEmailVerifiedWhereNameActivationKey(ctx, true, ua.userName, ua.activationKey)
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
