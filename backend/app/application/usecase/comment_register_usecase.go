package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type RegisterCommentUseCaseInterface interface {
	RegisterCommentUseCase() (*model.Comment, error)
}

type RegisterComment struct {
	ctx                context.Context
	tx                 mysql.DBTransaction
	tokenUserName      string
	postingID          int
	reqRegisterComment *modelHTTP.Comment
	userRepo           *repository.UserRepository
	postingRepo        *repository.PostingRepository
	commentRepo        *repository.CommentRepository
	notificationRepo   *repository.NotificationRepository
}

func NewRegisterComment(ctx context.Context, tx mysql.DBTransaction, tokenUserName string, postingID int, reqRegisterComment *modelHTTP.Comment, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository, commentRepo *repository.CommentRepository, notificationRepo *repository.NotificationRepository) *RegisterComment {
	return &RegisterComment{
		ctx:                ctx,
		tx:                 tx,
		tokenUserName:      tokenUserName,
		postingID:          postingID,
		reqRegisterComment: reqRegisterComment,
		userRepo:           userRepo,
		postingRepo:        postingRepo,
		commentRepo:        commentRepo,
		notificationRepo:   notificationRepo,
	}
}

func (comment *RegisterComment) RegisterCommentUseCase() error {
	// check userName in token exists
	_, err := comment.userRepo.GetUserWhereName(comment.ctx, comment.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return lib.ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	_, err = comment.postingRepo.GetWhereID(comment.ctx, int64(comment.postingID))
	if err != nil {
		return err
	}
	err = comment.tx.Do(comment.ctx, func(ctx context.Context) error {
		c := model.Comment{
			UserName:  comment.tokenUserName,
			PostingID: int64(comment.postingID),
			Comment:   comment.reqRegisterComment.Comment,
		}
		err := comment.commentRepo.Create(ctx, &c)
		if err != nil {
			return err
		}

		// TODO notification
		// if comment.userName != p.tokenUserName {
		// 	n := model.Notification{
		// 		VisitorName: comment.userName,
		// 		VisitedName: p.tokenUserName,
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
