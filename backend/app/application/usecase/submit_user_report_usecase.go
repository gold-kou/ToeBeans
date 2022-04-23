package usecase

import (
	"context"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	modelHTTP "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
)

type SubmitUserReportUseCaseInterface interface {
	SubmitUserReportUseCase(context.Context) error
}

type SubmitUserReport struct {
	tx                  mysql.DBTransaction
	userName            string
	reqSubmitUserReport *modelHTTP.RequestSubmitUserReport
	userRepo            *repository.UserRepository
	userReportRepo      *repository.UserReportRepository
}

func NewSubmitUserReport(tx mysql.DBTransaction, userName string, reqSubmitUserReport *modelHTTP.RequestSubmitUserReport, userRepo *repository.UserRepository, userReportRepo *repository.UserReportRepository) *SubmitUserReport {
	return &SubmitUserReport{
		tx:                  tx,
		userName:            userName,
		reqSubmitUserReport: reqSubmitUserReport,
		userRepo:            userRepo,
		userReportRepo:      userReportRepo,
	}
}

func (ur *SubmitUserReport) SubmitUserReportUseCase(ctx context.Context) error {
	_, err := ur.userRepo.GetUserWhereName(ctx, ur.userName)
	if err != nil {
		if err == repository.ErrNotExistsData {
			return ErrNotExitsUser
		}
	}

	userReport := model.UserReport{UserName: ur.userName, Detail: ur.reqSubmitUserReport.Detail}
	err = ur.userReportRepo.Create(ctx, &userReport)
	if err != nil {
		return err
	}

	return nil
}
