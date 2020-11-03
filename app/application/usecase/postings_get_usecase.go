package usecase

import (
	"context"
	"time"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/model"
	"github.com/gold-kou/ToeBeans/app/domain/repository"
)

type GetPostingsUseCaseInterface interface {
	GetPostingsUseCase() (*model.Posting, error)
}

type GetPostings struct {
	ctx         context.Context
	tx          mysql.DBTransaction
	sinceAt     time.Time
	limit       int8
	userName    string
	postingRepo *repository.PostingRepository
}

func NewGetPostings(ctx context.Context, tx mysql.DBTransaction, sinceAt time.Time, limit int8, userName string, postingRepo *repository.PostingRepository) *GetPostings {
	return &GetPostings{
		ctx:         ctx,
		tx:          tx,
		sinceAt:     sinceAt,
		limit:       limit,
		userName:    userName,
		postingRepo: postingRepo,
	}
}

func (p *GetPostings) GetPostingsUseCase() (postings []model.Posting, err error) {
	postings, err = p.postingRepo.GetPostings(p.ctx, p.sinceAt, p.limit, p.userName)
	if err != nil {
		return
	}
	return
}
