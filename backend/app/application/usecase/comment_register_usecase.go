package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type RegisterCommentUseCaseInterface interface {
	RegisterCommentUseCase() (*model.Comment, error)
}

type RegisterComment struct {
	tx                 mysql.DBTransaction
	tokenUserID        int64
	tokenUserName      string
	postingID          int
	reqRegisterComment *modelHTTP.RequestRegisterComment
	userRepo           *repository.UserRepository
	postingRepo        *repository.PostingRepository
	commentRepo        *repository.CommentRepository
	notificationRepo   *repository.NotificationRepository
}

func NewRegisterComment(tx mysql.DBTransaction, tokenUserID int64, tokenUserName string, postingID int, reqRegisterComment *modelHTTP.RequestRegisterComment, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository, commentRepo *repository.CommentRepository, notificationRepo *repository.NotificationRepository) *RegisterComment {
	return &RegisterComment{
		tx:                 tx,
		tokenUserID:        tokenUserID,
		tokenUserName:      tokenUserName,
		postingID:          postingID,
		reqRegisterComment: reqRegisterComment,
		userRepo:           userRepo,
		postingRepo:        postingRepo,
		commentRepo:        commentRepo,
		notificationRepo:   notificationRepo,
	}
}

func (comment *RegisterComment) RegisterCommentUseCase(ctx context.Context) error {
	// check userName in token exists
	_, err := comment.userRepo.GetUserWhereName(ctx, comment.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrTokenInvalidNotExistingUserName
		}
		return err
	}

	_, err = comment.postingRepo.GetWhereID(ctx, int64(comment.postingID))
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrNotExistsData
		}
		return err
	}
	err = comment.tx.Do(ctx, func(ctx context.Context) error {
		c := model.Comment{
			UserID:    comment.tokenUserID,
			PostingID: int64(comment.postingID),
			Comment:   comment.reqRegisterComment.Comment,
		}
		err := comment.commentRepo.Create(ctx, &c)
		if err != nil {
			if err == repository.ErrDuplicateData {
				return ErrDuplicateData
			}
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
