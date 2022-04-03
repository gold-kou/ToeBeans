package repository

import (
	"context"
	"database/sql"
	"time"

	m "github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

type PasswordResetRepositoryInterface interface {
	FindByUserID(ctx context.Context, userID int64) (passwordReset model.PasswordReset, err error)
	UpdateWhereUserID(ctx context.Context, count uint16, resetKey string, expiresAt time.Time, userID int64) (err error)
	DeleteWhereUserID(ctx context.Context, userID int64) (err error)
}

type PasswordResetRepository struct {
	db *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{
		db: db,
	}
}

func (r *PasswordResetRepository) FindByUserID(ctx context.Context, userID int64) (passwordReset model.PasswordReset, err error) {
	q := "SELECT `id`, `user_id`, `password_reset_email_count`, `password_reset_key`, `password_reset_key_expires_at`, `created_at`, `updated_at` FROM `password_resets` WHERE `user_id` = ?"
	err = r.db.QueryRowContext(ctx, q, userID).Scan(&passwordReset.ID, &passwordReset.UserID, &passwordReset.PasswordResetEmailCount, &passwordReset.PasswordResetKey, &passwordReset.PasswordResetKeyExpiresAt, &passwordReset.CreatedAt, &passwordReset.UpdatedAt)
	if err == sql.ErrNoRows {
		err = ErrNotExistsData
		return
	}
	return
}

func (r *PasswordResetRepository) UpdateWhereUserID(ctx context.Context, count uint8, resetKey string, expiresAt time.Time, userID int64) (err error) {
	q := "UPDATE `password_resets` SET `password_reset_email_count` = ?, `password_reset_key` = ?, `password_reset_key_expires_at` = ? WHERE `user_id` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, count, resetKey, expiresAt, userID)
	} else {
		_, err = r.db.ExecContext(ctx, q, count, resetKey, expiresAt, userID)
	}
	return
}

func (r *PasswordResetRepository) DeleteWhereUserID(ctx context.Context, userID int64) (err error) {
	q := "DELETE FROM `password_resets` WHERE `user_id` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userID)
	} else {
		_, err = r.db.ExecContext(ctx, q, userID)
	}
	return
}
