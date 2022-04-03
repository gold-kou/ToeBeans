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
	tx                    mysql.DBTransaction
	reqPasswordResetEmail *modelHTTP.Email
	userRepo              *repository.UserRepository
	passwordResetRepo     *repository.PasswordResetRepository
}

func NewPasswordResetEmail(tx mysql.DBTransaction, reqPasswordResetEmail *modelHTTP.Email, userRepo *repository.UserRepository, passwordResetRepo *repository.PasswordResetRepository) *PasswordResetEmail {
	return &PasswordResetEmail{
		tx:                    tx,
		reqPasswordResetEmail: reqPasswordResetEmail,
		userRepo:              userRepo,
		passwordResetRepo:     passwordResetRepo,
	}
}

func (re *PasswordResetEmail) PasswordResetEmailUseCase(ctx context.Context) (err error) {
	// email exists check
	u, err := re.userRepo.GetUserWhereEmail(ctx, re.reqPasswordResetEmail.Email)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrNotExistsData
		}
		return
	}

	// count over check
	var count uint8
	pr, err := re.passwordResetRepo.FindByUserID(ctx, u.ID)
	if err != nil {
		if err == repository.ErrNotExistsData {
			count = 0
		} else {
			return err
		}
	}
	if pr.PasswordResetEmailCount != 0 {
		count = pr.PasswordResetEmailCount
	}
	if count >= model.MaxLimitPasswordResetPerDay {
		return ErrOverPasswordResetCount
	}

	// set password reset key
	resetKey, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	err = re.passwordResetRepo.UpdateWhereUserID(ctx, count+1, resetKey.String(), lib.NowFunc().Add(24*time.Hour), u.ID)
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
