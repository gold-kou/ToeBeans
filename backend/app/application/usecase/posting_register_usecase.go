package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/google/uuid"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/gcp"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/aws"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

var ErrNotCatImage = errors.New("you can post only a cat image")

type RegisterPostingUseCaseInterface interface {
	RegisterPostingUseCase() (*model.Posting, error)
}

type RegisterPosting struct {
	ctx                context.Context
	tx                 mysql.DBTransaction
	tokenUserName      string
	reqRegisterPosting *modelHTTP.RequestRegisterPosting
	userRepo           *repository.UserRepository
	postingRepo        *repository.PostingRepository
}

func NewRegisterPosting(ctx context.Context, tx mysql.DBTransaction, tokenUserName string, reqRegisterPosting *modelHTTP.RequestRegisterPosting, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository) *RegisterPosting {
	return &RegisterPosting{
		ctx:                ctx,
		tx:                 tx,
		tokenUserName:      tokenUserName,
		reqRegisterPosting: reqRegisterPosting,
		userRepo:           userRepo,
		postingRepo:        postingRepo,
	}
}

func (posting *RegisterPosting) RegisterPostingUseCase() error {
	// check userName in token exists
	_, err := posting.userRepo.GetUserWhereName(posting.ctx, posting.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return lib.ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	// base64 decode
	decodedImg, err := base64.StdEncoding.DecodeString(posting.reqRegisterPosting.Image)
	if err != nil {
		return ErrDecodeImage
	}

	// save decoded file
	u, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	file, err := os.Create("image" + u.String() + ".jpg")
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
		_ = os.Remove("image" + u.String() + ".jpg")
	}()

	// check cat or not
	if flag.Lookup("test.v") == nil {
		labels, err := gcp.DetectLabels("image" + u.String() + ".jpg")
		if err != nil {
			return err
		}
		for i, l := range labels {
			if l == "Cat" || strings.Contains(l, "cat") {
				break
			}
			if i == len(labels)-1 {
				return ErrNotCatImage
			}
		}
	}

	// put decoded file to s3
	key := lib.NowFunc().Format(lib.DateTimeFormatNoSeparator) + "_" + posting.tokenUserName
	o, err := aws.UploadObject(os.Getenv("S3_BUCKET_POSTINGS"), key, decodedImg)
	if err != nil {
		return err
	}

	// INSERT
	if os.Getenv("APP_ENV") == "development" {
		// localhostに置換したが、ブラウザの仕様でCORBされる。imgタグでオリジン跨ぎの画像を読み込みできない。
		// with MIME type text/html. See https://www.chromestatus.com/feature/5629709824032768 for more details.
		o.Location = strings.Replace(o.Location, "minio", "localhost", 1)
	}
	err = posting.tx.Do(posting.ctx, func(ctx context.Context) error {
		u := model.Posting{
			UserName: posting.tokenUserName,
			Title:    posting.reqRegisterPosting.Title,
			ImageURL: o.Location,
		}
		err = posting.postingRepo.Create(ctx, &u)
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
