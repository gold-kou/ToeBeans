package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/app/lib"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/model"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

type GetCommentsUseCaseInterface interface {
	GetCommentsUseCase() (*model.Comment, error)
}

type GetComments struct {
	ctx           context.Context
	tx            mysql.DBTransaction
	tokenUserName string
	postingID     int64
	userRepo      *repository.UserRepository
	postingRepo   *repository.PostingRepository
	commentRepo   *repository.CommentRepository
}

func NewGetComments(ctx context.Context, tx mysql.DBTransaction, tokenUserName string, postingID int64, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository, commentRepo *repository.CommentRepository) *GetComments {
	return &GetComments{
		ctx:           ctx,
		tx:            tx,
		tokenUserName: tokenUserName,
		postingID:     postingID,
		userRepo:      userRepo,
		postingRepo:   postingRepo,
		commentRepo:   commentRepo,
	}
}

func (c *GetComments) GetCommentsUseCase() (comments []model.Comment, err error) {
	// check userName in token exists
	_, err = c.userRepo.GetUserWhereName(c.ctx, c.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return nil, lib.ErrTokenInvalidNotExistingUserName
		}
		return nil, err
	}

	// check postingID exists
	_, err = c.postingRepo.GetWhereID(c.ctx, c.postingID)
	if err != nil {
		return
	}

	comments, err = c.commentRepo.GetCommentsWherePostingID(c.ctx, c.postingID)
	if err != nil {
		if err == repository.ErrNotExistsData {
			// not error
			return nil, nil
		}
	}
	return
}
