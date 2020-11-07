package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/model"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

type DeleteCommentUseCaseInterface interface {
	DeleteCommentUseCase() (*model.Comment, error)
}

type DeleteComment struct {
	ctx         context.Context
	tx          mysql.DBTransaction
	userName    string
	commentID   int64
	commentRepo *repository.CommentRepository
}

func NewDeleteComment(ctx context.Context, tx mysql.DBTransaction, userName string, commentID int64, commentRepo *repository.CommentRepository) *DeleteComment {
	return &DeleteComment{
		ctx:         ctx,
		tx:          tx,
		userName:    userName,
		commentID:   commentID,
		commentRepo: commentRepo,
	}
}

func (comment *DeleteComment) DeleteCommentUseCase() error {
	if err := comment.commentRepo.DeleteWhereID(comment.ctx, comment.commentID); err != nil {
		return err
	}
	return nil
}
