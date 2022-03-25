package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type GetCommentsUseCaseInterface interface {
	GetCommentsUseCase() (*model.Comment, error)
}

type GetComments struct {
	tx            mysql.DBTransaction
	tokenUserName string
	postingID     int64
	userRepo      *repository.UserRepository
	postingRepo   *repository.PostingRepository
	commentRepo   *repository.CommentRepository
}

func NewGetComments(tx mysql.DBTransaction, tokenUserName string, postingID int64, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository, commentRepo *repository.CommentRepository) *GetComments {
	return &GetComments{
		tx:            tx,
		tokenUserName: tokenUserName,
		postingID:     postingID,
		userRepo:      userRepo,
		postingRepo:   postingRepo,
		commentRepo:   commentRepo,
	}
}

func (c *GetComments) GetCommentsUseCase(ctx context.Context) (comments []model.Comment, err error) {
	// check userName in token exists
	_, err = c.userRepo.GetUserWhereName(ctx, c.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return nil, ErrTokenInvalidNotExistingUserName
		}
		return nil, err
	}

	// check postingID exists
	_, err = c.postingRepo.GetWhereID(ctx, c.postingID)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return nil, ErrNotExistsData
		}
		return
	}

	comments, err = c.commentRepo.GetCommentsWherePostingID(ctx, c.postingID)
	if err != nil {
		if err == repository.ErrNotExistsData {
			// not error
			return nil, nil
		}
	}
	return
}
