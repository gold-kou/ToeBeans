package usecase

import (
	"context"
	"errors"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrNotVerifiedUser = errors.New("not email verified user")

type LoginUseCaseInterface interface {
	LoginUseCase() (string, error)
}

type Login struct {
	tx       mysql.DBTransaction
	reqLogin *modelHTTP.RequestLogin
	userRepo *repository.UserRepository
}

func NewLogin(tx mysql.DBTransaction, reqLogin *modelHTTP.RequestLogin, userRepo *repository.UserRepository) *Login {
	return &Login{
		tx:       tx,
		reqLogin: reqLogin,
		userRepo: userRepo,
	}
}

func (l *Login) LoginUseCase(ctx context.Context) (idToken string, err error) {
	user, err := l.userRepo.GetUserWhereEmail(ctx, l.reqLogin.Email)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return "", ErrNotExistsData
		}
		return
	}

	if !user.EmailVerified {
		return "", ErrNotVerifiedUser
	}

	// password check
	if user.Name != helper.GuestUserName {
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(l.reqLogin.Password)); err != nil {
			return "", ErrNotCorrectPassword
		}
	}

	// generate token
	idToken, err = helper.GenerateToken(user.ID, user.Name)
	if err != nil {
		return
	}

	return
}
