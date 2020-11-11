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
	postingRepo        *repository.PostingRepository
	commentRepo        *repository.CommentRepository
	notificationRepo   *repository.NotificationRepository
}

func NewRegisterComment(ctx context.Context, tx mysql.DBTransaction, userName string, reqRegisterComment *modelHTTP.Comment, postingRepo *repository.PostingRepository, commentRepo *repository.CommentRepository, notificationRepo *repository.NotificationRepository) *RegisterComment {
	return &RegisterComment{
		ctx:                ctx,
		tx:                 tx,
		userName:           userName,
		reqRegisterComment: reqRegisterComment,
		postingRepo:        postingRepo,
		commentRepo:        commentRepo,
		notificationRepo:   notificationRepo,
	}
}

func (comment *RegisterComment) RegisterCommentUseCase() error {
	_, err := comment.postingRepo.GetWhereID(comment.ctx, comment.reqRegisterComment.PostingId)
	if err != nil {
		return err
	}
	err = comment.tx.Do(comment.ctx, func(ctx context.Context) error {
		c := model.Comment{
			UserName:  comment.userName,
			PostingID: comment.reqRegisterComment.PostingId,
			Comment:   comment.reqRegisterComment.Comment,
		}
		_, err := comment.commentRepo.Create(ctx, &c)
		if err != nil {
			return err
		}

		// TODO notification
		// if comment.userName != p.UserName {
		// 	n := model.Notification{
		// 		VisitorName: comment.userName,
		// 		VisitedName: p.UserName,
		// 		Action:      model.CommentAction,
		// 	}
		// 	if err = comment.notificationRepo.Create(ctx, &n); err != nil {
		// 		return err
		// 	}
		// }

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
