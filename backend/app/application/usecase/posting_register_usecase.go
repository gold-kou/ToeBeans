package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/google/uuid"

	"github.com/gold-kou/ToeBeans/backend/app"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/aws"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/gcp"
	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
	"github.com/gold-kou/ToeBeans/backend/app/lib"
)

var ErrNotCatImage = errors.New("you can post only a cat image")
var bucketPosting string

func init() {
	bucketPosting = os.Getenv("S3_BUCKET_POSTINGS")
	if bucketPosting == "" {
		panic("S3_BUCKET_POSTINGS is unset")
	}
}

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
			return ErrTokenInvalidNotExistingUserName
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

	// check cat or not
	if flag.Lookup("test.v") == nil {
		labels, err := gcp.DetectLabels(filePath)
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

	// put file to s3
	savedFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer savedFile.Close()

	key := lib.NowFunc().Format(lib.DateTimeFormatNoSeparator) + "_" + posting.tokenUserName
	o, err := aws.UploadObject(bucketPosting, key, savedFile)
	if err != nil {
		return err
	}

	// INSERT
	if app.IsLocal() {
		o.Location = strings.Replace(o.Location, "minio", "localhost", 1)
	}
	err = posting.tx.Do(posting.ctx, func(ctx context.Context) error {
		p := model.Posting{
			UserName: posting.tokenUserName,
			Title:    posting.reqRegisterPosting.Title,
			ImageURL: o.Location,
			// ImageURL: "http://localhost:9000/toebeans-postings/20200101000000_testUser1",
		}
		err = posting.postingRepo.Create(ctx, &p)
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
