package usecase

import (
	"context"
	"time"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type GetPostingsUseCaseInterface interface {
	GetPostingsUseCase() (*model.Posting, error)
}

type GetPostings struct {
	tx             mysql.DBTransaction
	tokenUserName  string
	sinceAt        time.Time
	limit          int8
	targetUserName string
	userRepo       *repository.UserRepository
	postingRepo    *repository.PostingRepository
	likeRepo       *repository.LikeRepository
}

func NewGetPostings(tx mysql.DBTransaction, tokenUserName string, sinceAt time.Time, limit int8, targetUserName string, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository, likeRepo *repository.LikeRepository) *GetPostings {
	return &GetPostings{
		tx:             tx,
		tokenUserName:  tokenUserName,
		sinceAt:        sinceAt,
		limit:          limit,
		targetUserName: targetUserName,
		userRepo:       userRepo,
		postingRepo:    postingRepo,
		likeRepo:       likeRepo,
	}
}

func (p *GetPostings) GetPostingsUseCase(ctx context.Context) (postings []model.Posting, userNames []string, likedCounts []int64, likes []model.Like, err error) {
	// check userName in token exists
	tokenUser, err := p.userRepo.GetUserWhereName(ctx, p.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			err = ErrTokenInvalidNotExistingUserName
			return
		}
		return
	}

	// check targetUserName exists
	var targetUser model.User
	if p.targetUserName != "" {
		targetUser, err = p.userRepo.GetUserWhereName(ctx, p.targetUserName)
		if err != nil {
			if err == repository.ErrNotExistsData {
				err = ErrNotExistsData
				return
			}
			return
		}
	}

	likes, err = p.likeRepo.GetWhereUserID(ctx, tokenUser.ID)
	if err != nil {
		if err == repository.ErrNotExistsData {
			// not error
			err = nil
		}
		return
	}

	postings, err = p.postingRepo.GetPostings(ctx, p.sinceAt, p.limit, targetUser.ID)
	if err != nil {
		if err == repository.ErrNotExistsData {
			// not error
			err = nil
		}
		return
	}

	for _, posting := range postings {
		var user model.User
		user, err = p.userRepo.GetUserWhereID(ctx, posting.UserID)
		if err != nil {
			// ここでのnot exists errorは500エラー
			return
		}
		userNames = append(userNames, user.Name)

		var likedCount int64
		likedCount, err = p.likeRepo.GetLikedCountWherePostingID(ctx, posting.ID)
		if err != nil {
			return
		}
		likedCounts = append(likedCounts, likedCount)
	}
	return
}
