package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type DeleteCommentUseCaseInterface interface {
	DeleteCommentUseCase() (*model.Comment, error)
}

type DeleteComment struct {
	ctx           context.Context
	tx            mysql.DBTransaction
	tokenUserName string
	commentID     int64
	userRepo      *repository.UserRepository
	commentRepo   *repository.CommentRepository
}

func NewDeleteComment(ctx context.Context, tx mysql.DBTransaction, tokenUserName string, commentID int64, userRepo *repository.UserRepository, commentRepo *repository.CommentRepository) *DeleteComment {
	return &DeleteComment{
		ctx:           ctx,
		tx:            tx,
		tokenUserName: tokenUserName,
		commentID:     commentID,
		userRepo:      userRepo,
		commentRepo:   commentRepo,
	}
}

func (comment *DeleteComment) DeleteCommentUseCase() error {
	// check userName in token exists
	_, err := comment.userRepo.GetUserWhereName(comment.ctx, comment.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return lib.ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	if err := comment.commentRepo.DeleteWhereID(comment.ctx, comment.commentID); err != nil {
		return err
	}
	return nil
}
