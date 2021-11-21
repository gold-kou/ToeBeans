package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
	"github.com/gold-kou/ToeBeans/backend/app/lib"
	"golang.org/x/crypto/bcrypt"
)

type LoginUseCaseInterface interface {
	LoginUseCase() (string, error)
}

type Login struct {
	ctx      context.Context
	tx       mysql.DBTransaction
	reqLogin *modelHTTP.RequestLogin
	userRepo *repository.UserRepository
}

func NewLogin(ctx context.Context, tx mysql.DBTransaction, reqLogin *modelHTTP.RequestLogin, userRepo *repository.UserRepository) *Login {
	return &Login{
		ctx:      ctx,
		tx:       tx,
		reqLogin: reqLogin,
		userRepo: userRepo,
	}
}

func (l *Login) LoginUseCase() (idToken string, err error) {
	user, err := l.userRepo.GetUserWhereEmail(l.ctx, l.reqLogin.Email)
	if err != nil {
		return
	}

	if !user.EmailVerified {
		return "", ErrNotVerifiedUser
	}

	// password check
	if user.Name != lib.GuestUserName {
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(l.reqLogin.Password)); err != nil {
			return "", ErrNotCorrectPassword
		}
	}

	// generate token
	idToken, err = lib.GenerateToken(user.Name)
	if err != nil {
		return
	}

	return
}
