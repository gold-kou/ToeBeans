package usecase

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/model"
	modelHTTP "github.com/gold-kou/ToeBeans/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

type RegisterUserUseCaseInterface interface {
	RegisterUserUseCase() (*model.User, error)
}

type RegisterUser struct {
	ctx             context.Context
	tx              mysql.DBTransaction
	reqRegisterUser *modelHTTP.RequestRegisterUser
	userRepo        *repository.UserRepository
}

func NewRegisterUser(ctx context.Context, tx mysql.DBTransaction, reqRegisterUser *modelHTTP.RequestRegisterUser, userRepo *repository.UserRepository) *RegisterUser {
	return &RegisterUser{
		ctx:             ctx,
		tx:              tx,
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

	// send an email
	//  title := "Welcome to ToeBeans"
	//  activateLink := fmt.Sprintf("https://<domain>/user-activation/%s/%s", user.reqRegisterUser.UserName, activationKey.String())
	//  body := fmt.Sprintf("Hi " +
	//	  user.reqRegisterUser.UserName +
	//	  ",\n" +
	//	  "\n" +
	//	  "You are just one step away from activating your account on the ToeBeans! Click on the link and start enjoying your account:\n" +
	//	  activateLink +
	//	  "\n" +
	//	  "\n" +
	//	  "Didn't ask for a new account? You can ignore this email.")
	//  err = aws.SendEmail(user.reqRegisterUser.Email, title, body)
	//  if err != nil {
	//	  return err
	//  }

	err = user.tx.Do(user.ctx, func(ctx context.Context) error {
		u := model.User{
			Name:          user.reqRegisterUser.UserName,
			Email:         user.reqRegisterUser.Email,
			Password:      user.reqRegisterUser.Password,
			ActivationKey: activationKey.String(),
		}
		err = user.userRepo.Create(ctx, &u)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
