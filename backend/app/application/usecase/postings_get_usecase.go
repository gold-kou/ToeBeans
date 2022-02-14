package usecase

import (
	"context"
	"time"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
	"github.com/gold-kou/ToeBeans/backend/app/lib"
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

func (p *GetPostings) GetPostingsUseCase() (postings []model.Posting, likes []model.Like, err error) {
	// check userName in token exists
	_, err = p.userRepo.GetUserWhereName(p.ctx, p.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return nil, nil, lib.ErrTokenInvalidNotExistingUserName
		}
		return nil, nil, err
	}

	// check userName exists
	if p.userName != "" {
		_, err = p.userRepo.GetUserWhereName(p.ctx, p.userName)
		if err != nil {
			return
		}
	}

	likes, err = p.likeRepo.GetWhereUserName(p.ctx, p.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			// not error
			return nil, nil, nil
		}
		return
	}

	postings, err = p.postingRepo.GetPostings(p.ctx, p.sinceAt, p.limit, p.userName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			// not error
			return nil, nil, nil
		}
		return
	}
	return
}
