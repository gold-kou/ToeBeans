package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type DeleteCommentUseCaseInterface interface {
	DeleteCommentUseCase() (*model.Comment, error)
}

type DeleteComment struct {
	tx            mysql.DBTransaction
	tokenUserName string
	commentID     int64
	userRepo      *repository.UserRepository
	commentRepo   *repository.CommentRepository
}

func NewDeleteComment(tx mysql.DBTransaction, tokenUserName string, commentID int64, userRepo *repository.UserRepository, commentRepo *repository.CommentRepository) *DeleteComment {
	return &DeleteComment{
		tx:            tx,
		tokenUserName: tokenUserName,
		commentID:     commentID,
		userRepo:      userRepo,
		commentRepo:   commentRepo,
	}
}

func (comment *DeleteComment) DeleteCommentUseCase(ctx context.Context) error {
	// check userName in token exists
	_, err := comment.userRepo.GetUserWhereName(ctx, comment.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	if err := comment.commentRepo.DeleteWhereID(ctx, comment.commentID); err != nil {
		if err == repository.ErrNotExistsData {
			return ErrNotExistsData
		}
		return err
	}
	return nil
}
