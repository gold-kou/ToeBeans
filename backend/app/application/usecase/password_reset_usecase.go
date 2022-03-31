package usecase

import (
	"context"
	"errors"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/lib"
	"golang.org/x/crypto/bcrypt"

	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

var ErrPasswordResetKeyWrong = errors.New("password reset key is wrong")
var ErrPasswordResetKeyExpired = errors.New("password reset key is expired")

type PasswordResetUseCaseInterface interface {
	PasswordResetUseCase() error
}

type PasswordReset struct {
	tx                mysql.DBTransaction
	reqPasswordReset  *modelHTTP.RequestResetPassword
	userRepo          *repository.UserRepository
	passwordResetRepo *repository.PasswordResetRepository
}

func NewPasswordReset(tx mysql.DBTransaction, reqPasswordReset *modelHTTP.RequestResetPassword, userRepo *repository.UserRepository, passwordResetRepo *repository.PasswordResetRepository) *PasswordReset {
	return &PasswordReset{
		tx:                tx,
		reqPasswordReset:  reqPasswordReset,
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
	}
}

func (reset *PasswordReset) PasswordResetUseCase(ctx context.Context) (err error) {
	// user name exists check
	u, err := reset.userRepo.GetUserWhereName(ctx, reset.reqPasswordReset.UserName)
	if err == repository.ErrNotExistsData {
		return ErrNotExistsData
	}

	// password reset check
	pr, err := reset.passwordResetRepo.FindByUserID(ctx, u.ID)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrNotExistsData
		}
		return err
	}
	if pr.PasswordResetKey != reset.reqPasswordReset.PasswordResetKey {
		return ErrPasswordResetKeyWrong
	}
	if pr.PasswordResetKeyExpiresAt.Before(lib.NowFunc()) {
		return ErrPasswordResetKeyExpired
	}

	// reset password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reset.reqPasswordReset.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := reset.userRepo.ResetPassword(ctx, string(hashedPassword), reset.reqPasswordReset.UserName); err != nil {
		return err
	}
	return
}
