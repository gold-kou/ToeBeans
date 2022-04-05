package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type SubmitPostingReportUseCaseInterface interface {
	SubmitPostingReportUseCase(context.Context) error
}

type SubmitPostingReport struct {
	tx                     mysql.DBTransaction
	postingID              int
	reqSubmitPostingReport *modelHTTP.RequestSubmitPostingReport
	postingRepo            *repository.PostingRepository
	postingReportRepo      *repository.PostingReportRepository
}

func NewSubmitPostingReport(tx mysql.DBTransaction, postingID int, reqSubmitPostingReport *modelHTTP.RequestSubmitPostingReport, postingRepo *repository.PostingRepository, postingReportRepo *repository.PostingReportRepository) *SubmitPostingReport {
	return &SubmitPostingReport{
		tx:                     tx,
		postingID:              postingID,
		reqSubmitPostingReport: reqSubmitPostingReport,
		postingRepo:            postingRepo,
		postingReportRepo:      postingReportRepo,
	}
}

func (pr *SubmitPostingReport) SubmitPostingReportUseCase(ctx context.Context) error {
	_, err := pr.postingRepo.GetWhereID(ctx, int64(pr.postingID))
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrNotExistsData
		}
	}

	postingReport := model.PostingReport{PostingID: pr.postingID, Detail: pr.reqSubmitPostingReport.Detail}
	err = pr.postingReportRepo.Create(ctx, &postingReport)
	if err != nil {
		return err
	}

	return nil
}
