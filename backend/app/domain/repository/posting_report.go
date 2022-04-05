package repository

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
	m "github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

type PostingReportRepositoryInterface interface {
	Create(ctx context.Context, user *model.PostingReport) (err error)
}

type PostingReportRepository struct {
	db *sql.DB
}

func NewPostingReportRepository(db *sql.DB) *PostingReportRepository {
	return &PostingReportRepository{
		db: db,
	}
}

func (r *PostingReportRepository) Create(ctx context.Context, postingReport *model.PostingReport) (err error) {
	q := "INSERT INTO `posting_reports` (`posting_id`, `detail`) VALUES (?, ?)"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, postingReport.PostingID, postingReport.Detail)
	} else {
		_, err = r.db.ExecContext(ctx, q, postingReport.PostingID, postingReport.Detail)
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok && mysqlErr.Number == 1062 {
		return ErrDuplicateData
	}
	return
}
