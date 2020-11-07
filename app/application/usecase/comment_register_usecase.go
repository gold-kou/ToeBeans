package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/model"
	modelHTTP "github.com/gold-kou/ToeBeans/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

type RegisterCommentUseCaseInterface interface {
	RegisterCommentUseCase() (*model.Comment, error)
}

type RegisterComment struct {
	ctx                context.Context
	tx                 mysql.DBTransaction
	userName           string
	reqRegisterComment *modelHTTP.Comment
	commentRepo        *repository.CommentRepository
}

func NewRegisterComment(ctx context.Context, tx mysql.DBTransaction, userName string, reqRegisterComment *modelHTTP.Comment, commentRepo *repository.CommentRepository) *RegisterComment {
	return &RegisterComment{
		ctx:                ctx,
		tx:                 tx,
		userName:           userName,
		reqRegisterComment: reqRegisterComment,
		commentRepo:        commentRepo,
	}
}

func (comment *RegisterComment) RegisterCommentUseCase() error {
	c := model.Comment{
		UserName:  comment.userName,
		PostingID: comment.reqRegisterComment.PostingId,
		Comment:   comment.reqRegisterComment.Comment,
	}
	if err := comment.commentRepo.Create(comment.ctx, &c); err != nil {
		return err
	}
	return nil
}
