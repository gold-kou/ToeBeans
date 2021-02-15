package usecase

import (
	"context"
	"time"

	"github.com/gold-kou/ToeBeans/app/lib"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/model"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
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
}

func NewGetPostings(ctx context.Context, tx mysql.DBTransaction, tokenUserName string, sinceAt time.Time, limit int8, userName string, userRepo *repository.UserRepository, postingRepo *repository.PostingRepository) *GetPostings {
	return &GetPostings{
		ctx:           ctx,
		tx:            tx,
		tokenUserName: tokenUserName,
		sinceAt:       sinceAt,
		limit:         limit,
		userName:      userName,
		userRepo:      userRepo,
		postingRepo:   postingRepo,
	}
}

func (p *GetPostings) GetPostingsUseCase() (postings []model.Posting, err error) {
	// check userName in token exists
	_, err = p.userRepo.GetUserWhereName(p.ctx, p.tokenUserName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return nil, lib.ErrTokenInvalidNotExistingUserName
		}
		return nil, err
	}

	// check userName exists
	if p.userName != "" {
		_, err = p.userRepo.GetUserWhereName(p.ctx, p.userName)
		if err != nil {
			return
		}
	}

	postings, err = p.postingRepo.GetPostings(p.ctx, p.sinceAt, p.limit, p.userName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			// not error
			return nil, nil
		}
		return
	}
	return
}
