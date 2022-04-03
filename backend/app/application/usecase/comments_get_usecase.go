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

func (c *GetComments) GetCommentsUseCase(ctx context.Context) (comments []model.Comment, userNames []string, err error) {
	// check userName in token exists
	_, err = c.userRepo.GetUserWhereName(ctx, c.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			err = ErrTokenInvalidNotExistingUserName
			return
		}
		return
	}

	// check postingID exists
	_, err = c.postingRepo.GetWhereID(ctx, c.postingID)
	if err != nil {
		if err == repository.ErrNotExistsData {
			err = ErrNotExistsData
			return
		}
		return
	}

	comments, err = c.commentRepo.GetCommentsWherePostingID(ctx, c.postingID)
	if err != nil {
		if err == repository.ErrNotExistsData {
			// not error
			err = nil
			return
		}
	}
	for _, comment := range comments {
		var user model.User
		user, err = c.userRepo.GetUserWhereID(ctx, comment.UserID)
		if err != nil {
			// ここでのnot exists errorは500エラー
			return
		}
		userNames = append(userNames, user.Name)
	}
	return
}
