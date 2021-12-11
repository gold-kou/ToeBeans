package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/gold-kou/ToeBeans/backend/app/domain/model"

	"github.com/gold-kou/ToeBeans/backend/app/lib"
	"github.com/google/uuid"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"

	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

var ErrOverPasswordResetCount = errors.New("you can't reset password as it exceeds limit counts")

type PasswordResetEmailUseCaseInterface interface {
	PasswordResetEmailUseCase() error
}

type PasswordResetEmail struct {
	ctx                   context.Context
	tx                    mysql.DBTransaction
	reqPasswordResetEmail *modelHTTP.Email
	userRepo              *repository.UserRepository
}

func NewPasswordResetEmail(ctx context.Context, tx mysql.DBTransaction, reqPasswordResetEmail *modelHTTP.Email, userRepo *repository.UserRepository) *PasswordResetEmail {
	return &PasswordResetEmail{
		ctx:                   ctx,
		tx:                    tx,
		reqPasswordResetEmail: reqPasswordResetEmail,
		userRepo:              userRepo,
	}
}

func (re *PasswordResetEmail) PasswordResetEmailUseCase() (err error) {
	// email exists check
	u, err := re.userRepo.GetUserWhereEmail(re.ctx, re.reqPasswordResetEmail.Email)
	if err != nil {
		return
	}

	if u.PasswordResetEmailCount >= model.MaxLimitPasswordResetPerDay {
		return ErrOverPasswordResetCount
	}

	// set password reset key
	resetKey, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	err = re.userRepo.UpdatePasswordResetWhereEmail(re.ctx, u.PasswordResetEmailCount+1, resetKey.String(), lib.NowFunc().Add(24*time.Hour), re.reqPasswordResetEmail.Email)
	if err != nil {
		return
	}

	// send password reset key via an email
	//  title := "Reset your password on ToeBeans"
	//  resetLink := fmt.Sprintf("https://<domain>/reset-page/%s/%s", u.Name, resetKey)
	//  body := fmt.Sprintf("We got your request to change your password.\n" +
	//	  "\n" +
	//	  resetLink +
	//	  "\n" +
	//	  "\n" +
	//	  "Just so you know: You have 24 hours to pick your password. After that, you'll have to ask for a new one.\n" +
	//	  "Didn't ask for a new password? You can ignore this email.")
	//  err = aws.SendEmail(re.reqPasswordResetEmail.Email, title, body)
	//  if err != nil {
	//	  return
	//  }

	return
}
