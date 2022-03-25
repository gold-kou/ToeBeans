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
	tx               mysql.DBTransaction
	reqPasswordReset *modelHTTP.RequestResetPassword
	userRepo         *repository.UserRepository
}

func NewPasswordReset(tx mysql.DBTransaction, reqPasswordReset *modelHTTP.RequestResetPassword, userRepo *repository.UserRepository) *PasswordReset {
	return &PasswordReset{
		tx:               tx,
		reqPasswordReset: reqPasswordReset,
		userRepo:         userRepo,
	}
}

func (reset *PasswordReset) PasswordResetUseCase(ctx context.Context) (err error) {
	// name and resetKey exists check
	_, err = reset.userRepo.GetUserWhereNameResetKey(ctx, reset.reqPasswordReset.UserName, reset.reqPasswordReset.PasswordResetKey, lib.NowFunc())
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
	if err := reset.userRepo.ResetPassword(ctx, string(hashedPassword), reset.reqPasswordReset.UserName, reset.reqPasswordReset.PasswordResetKey); err != nil {
		return err
	}
	return
}
