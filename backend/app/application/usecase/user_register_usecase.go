package usecase

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/gold-kou/ToeBeans/backend/app"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/aws"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

var domain string

func init() {
	// testでの空文字を許容する
	domain = os.Getenv("DOMAIN")
}

type RegisterUserUseCaseInterface interface {
	RegisterUserUseCase() (*model.User, error)
}

type RegisterUser struct {
	ctx             context.Context
	tx              mysql.DBTransaction
	userName        string
	reqRegisterUser *modelHTTP.RequestRegisterUser
	userRepo        *repository.UserRepository
}

func NewRegisterUser(ctx context.Context, tx mysql.DBTransaction, userName string, reqRegisterUser *modelHTTP.RequestRegisterUser, userRepo *repository.UserRepository) *RegisterUser {
	return &RegisterUser{
		ctx:             ctx,
		tx:              tx,
		userName:        userName,
		reqRegisterUser: reqRegisterUser,
		userRepo:        userRepo,
	}
}

func (user *RegisterUser) RegisterUserUseCase() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.reqRegisterUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.reqRegisterUser.Password = string(hashedPassword)

	activationKey, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	err = user.tx.Do(user.ctx, func(ctx context.Context) error {
		u := model.User{
			Name:          user.userName,
			Email:         user.reqRegisterUser.Email,
			Password:      user.reqRegisterUser.Password,
			ActivationKey: activationKey.String(),
		}
		err = user.userRepo.Create(ctx, &u)
		if err != nil {
			if err == repository.ErrDuplicateData {
				return ErrDuplicateData
			}
			return err
		}

		// send an email
		if flag.Lookup("test.v") == nil {
			var prefix string
			if app.IsLocal() {
				prefix = "http://" + domain
			} else {
				prefix = "https://" + domain
			}
			title := "Welcome to ToeBeans"
			activateLink := fmt.Sprintf(prefix+"/user-activation/%s/%s", user.userName, activationKey.String())
			body := fmt.Sprintf("Hi " +
				user.userName +
				",\n" +
				"\n" +
				"You are just one step away from activating your account on the ToeBeans!" +
				"\n" +
				"Click on the link and start enjoying your account:\n" +
				activateLink +
				"\n" +
				"\n" +
				"Didn't ask for a new account? You can ignore this email.")
			err = aws.SendEmail(user.reqRegisterUser.Email, title, body)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
