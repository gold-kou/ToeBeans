package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/model"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

type GetCommentsUseCaseInterface interface {
	GetCommentsUseCase() (*model.Comment, error)
}

type GetComments struct {
	ctx         context.Context
	tx          mysql.DBTransaction
	postingID   int64
	commentRepo *repository.CommentRepository
}

func NewGetComments(ctx context.Context, tx mysql.DBTransaction, postingID int64, commentRepo *repository.CommentRepository) *GetComments {
	return &GetComments{
		ctx:         ctx,
		tx:          tx,
		postingID:   postingID,
		commentRepo: commentRepo,
	}
}

func (c *GetComments) GetCommentsUseCase() (comments []model.Comment, err error) {
	comments, err = c.commentRepo.GetCommentsWherePostingID(c.ctx, c.postingID)
	if err != nil {
		return
	}
	return
}
