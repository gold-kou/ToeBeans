package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"golang.org/x/crypto/bcrypt"

	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type PasswordResetUseCaseInterface interface {
	PasswordResetUseCase() error
}

type PasswordReset struct {
	ctx              context.Context
	tx               mysql.DBTransaction
	reqPasswordReset *modelHTTP.RequestResetPassword
	userRepo         *repository.UserRepository
}

func NewPasswordReset(ctx context.Context, tx mysql.DBTransaction, reqPasswordReset *modelHTTP.RequestResetPassword, userRepo *repository.UserRepository) *PasswordReset {
	return &PasswordReset{
		ctx:              ctx,
		tx:               tx,
		reqPasswordReset: reqPasswordReset,
		userRepo:         userRepo,
	}
}

func (reset *PasswordReset) PasswordResetUseCase() (err error) {
	// name and resetKey exists check
	_, err = reset.userRepo.GetUserWhereNameResetKey(reset.ctx, reset.reqPasswordReset.UserName, reset.reqPasswordReset.PasswordResetKey, lib.NowFunc())
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrNotExistsData
		}
		return err
	}

	// reset password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reset.reqPasswordReset.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := reset.userRepo.ResetPassword(reset.ctx, string(hashedPassword), reset.reqPasswordReset.UserName, reset.reqPasswordReset.PasswordResetKey); err != nil {
		return err
	}
	return
}
