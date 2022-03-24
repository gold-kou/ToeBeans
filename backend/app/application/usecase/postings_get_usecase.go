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
	ctx           context.Context
	tx            mysql.DBTransaction
	tokenUserName string
	sinceAt       time.Time
	limit         int8
	userName      string
	userRepo      *repository.UserRepository
	postingRepo   *repository.PostingRepository
	likeRepo      *repository.LikeRepository
}

func NewGetPostings(ctx context.Context, tx mysql.DBTransaction, tokenUserName string, sinceAt time.Time, limit int8, userName string, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository, likeRepo *repository.LikeRepository) *GetPostings {
	return &GetPostings{
		ctx:           ctx,
		tx:            tx,
		tokenUserName: tokenUserName,
		sinceAt:       sinceAt,
		limit:         limit,
		userName:      userName,
		userRepo:      userRepo,
		postingRepo:   postingRepo,
		likeRepo:      likeRepo,
	}
}

func (p *GetPostings) GetPostingsUseCase() (postings []model.Posting, likedCounts []int64, likes []model.Like, err error) {
	// check userName in token exists
	_, err = p.userRepo.GetUserWhereName(p.ctx, p.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			err = ErrTokenInvalidNotExistingUserName
			return
		}
		return
	}

	// check userName exists
	if p.userName != "" {
		_, err = p.userRepo.GetUserWhereName(p.ctx, p.userName)
		if err != nil {
			if err == repository.ErrNotExistsData {
				err = ErrNotExistsData
				return
			}
			return
		}
	}

	likes, err = p.likeRepo.GetWhereUserName(p.ctx, p.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			// not error
			err = nil
		}
		return
	}

	postings, err = p.postingRepo.GetPostings(p.ctx, p.sinceAt, p.limit, p.userName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			// not error
			err = nil
		}
		return
	}

	var likedCount int64
	for _, posting := range postings {
		likedCount, err = p.likeRepo.GetLikedCountWherePostingID(p.ctx, posting.ID)
		if err != nil {
			return
		}
		likedCounts = append(likedCounts, likedCount)
	}
	return
}
