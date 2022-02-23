package usecase

import (
	"context"
	"encoding/base64"
	"os"
	"strings"

	"github.com/gold-kou/ToeBeans/backend/app"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/aws"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

var bucketIcons string

func init() {
	bucketIcons = os.Getenv("S3_BUCKET_ICONS")
	if bucketIcons == "" {
		panic("S3_BUCKET_ICONS is unset")
	}
}

type UpdateUserUseCaseInterface interface {
	UpdateUserUseCase() (*model.User, error)
}

type UpdateUser struct {
	ctx           context.Context
	tx            mysql.DBTransaction
	userName      string
	reqUpdateUser *modelHTTP.RequestUpdateUser
	userRepo      *repository.UserRepository
}

func NewUpdateUser(ctx context.Context, tx mysql.DBTransaction, userName string, reqUpdateUser *modelHTTP.RequestUpdateUser, userRepo *repository.UserRepository) *UpdateUser {
	return &UpdateUser{
		ctx:           ctx,
		tx:            tx,
		userName:      userName,
		reqUpdateUser: reqUpdateUser,
		userRepo:      userRepo,
	}
}

func (user *UpdateUser) UpdateUserUseCase() error {
	// check user exists
	_, err := user.userRepo.GetUserWhereName(user.ctx, user.userName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrNotExitsUser
		}
		return err
	}

	// the case of password
	if user.reqUpdateUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.reqUpdateUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		err = user.userRepo.UpdatePasswordWhereName(user.ctx, string(hashedPassword), user.userName)
		if err != nil {
			return err
		}
	}
	// the case of icon
	if user.reqUpdateUser.Icon != "" {
		// base64 decode
		decodedImg, err := base64.StdEncoding.DecodeString(user.reqUpdateUser.Icon)
		if err != nil {
			return ErrDecodeImage
		}

		// save decoded file
		u, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		filePath := "image" + u.String() + ".jpg"
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.Write(decodedImg); err != nil {
			return err
		}
		if err := file.Sync(); err != nil {
			return err
		}
		// delete file
		defer func() {
			_ = os.Remove(filePath)
		}()

		savedFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer savedFile.Close()
		o, err := aws.UploadObject(bucketIcons, user.userName, savedFile)
		if err != nil {
			return err
		}

		if app.IsLocal() {
			o.Location = strings.Replace(o.Location, "minio", "localhost", 1)
		}
		err = user.userRepo.UpdateIconWhereName(user.ctx, o.Location, user.userName)
		if err != nil {
			return err
		}
	}
	// the case of self introduction
	if user.reqUpdateUser.SelfIntroduction != "" {
		err := user.userRepo.UpdateSelfIntroductionWhereName(user.ctx, user.reqUpdateUser.SelfIntroduction, user.userName)
		if err != nil {
			return err
		}
	}
	return nil
}
