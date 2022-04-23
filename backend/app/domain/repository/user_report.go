package repository

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
	m "github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

type UserReportRepositoryInterface interface {
	Create(ctx context.Context, user *model.UserReport) (err error)
}

type UserReportRepository struct {
	db *sql.DB
}

func NewUserReportRepository(db *sql.DB) *UserReportRepository {
	return &UserReportRepository{
		db: db,
	}
}

func (r *UserReportRepository) Create(ctx context.Context, userReport *model.UserReport) (err error) {
	q := "INSERT INTO `user_reports` (`user_name`, `detail`) VALUES (?, ?)"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userReport.UserName, userReport.Detail)
	} else {
		_, err = r.db.ExecContext(ctx, q, userReport.UserName, userReport.Detail)
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok && mysqlErr.Number == 1062 {
		return ErrDuplicateData
	}
	return
}
