package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type UpdateActivationUseCaseInterface interface {
	UserActivationUseCase() (*model.User, error)
}

type UserActivation struct {
	tx            mysql.DBTransaction
	userName      string
	activationKey string
	userRepo      *repository.UserRepository
}

func NewUserActivation(tx mysql.DBTransaction, userName, activationKey string, userRepo *repository.UserRepository) *UserActivation {
	return &UserActivation{
		tx:            tx,
		userName:      userName,
		activationKey: activationKey,
		userRepo:      userRepo,
	}
}

func (ua *UserActivation) UserActivationUseCase(ctx context.Context) error {
	err := ua.userRepo.UpdateEmailVerifiedWhereNameActivationKey(ctx, true, ua.userName, ua.activationKey)
	if err != nil {
		return err
	}
	return nil
}
