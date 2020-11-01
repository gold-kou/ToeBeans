package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/app/lib"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"golang.org/x/crypto/bcrypt"

	modelHTTP "github.com/gold-kou/ToeBeans/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
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

func (re *PasswordReset) PasswordResetUseCase() (err error) {
	// name and resetKey exists check
	_, err = re.userRepo.GetUserWhereNameResetKey(re.ctx, re.reqPasswordReset.UserName, re.reqPasswordReset.PasswordResetKey)
	if err != nil {
		return
	}

	// reset password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(re.reqPasswordReset.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err = re.userRepo.ResetPassword(re.ctx, string(hashedPassword), re.reqPasswordReset.UserName, re.reqPasswordReset.PasswordResetKey, lib.NowFunc()); err != nil {
		return
	}

	return
}
