package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/domain/model"

	"golang.org/x/crypto/bcrypt"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type ChangePasswordUseCaseInterface interface {
	ChangePasswordUseCase() (*model.User, error)
}

type ChangePassword struct {
	tx                mysql.DBTransaction
	tokenUserName     string
	reqChangePassword *modelHTTP.RequestChangePassword
	userRepo          *repository.UserRepository
}

func NewChangePassword(tx mysql.DBTransaction, tokenUserName string, reqChangePassword *modelHTTP.RequestChangePassword, userRepo *repository.UserRepository) *ChangePassword {
	return &ChangePassword{
		tx:                tx,
		tokenUserName:     tokenUserName,
		reqChangePassword: reqChangePassword,
		userRepo:          userRepo,
	}
}

func (user *ChangePassword) ChangePasswordUseCase(ctx context.Context) error {
	// check user exists
	dbUser, err := user.userRepo.GetUserWhereName(ctx, user.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	// old password check
	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.reqChangePassword.OldPassword)); err != nil {
		return ErrNotCorrectPassword
	}

	// change password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(user.reqChangePassword.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = user.userRepo.UpdatePasswordWhereName(ctx, string(hashedNewPassword), user.tokenUserName)
	if err != nil {
		return err
	}
	return nil
}
