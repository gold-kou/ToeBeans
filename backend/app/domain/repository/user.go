package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"

	m "github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *model.User) (err error)
	GetUserWhereEmail(ctx context.Context, email string) (user model.User, err error)
	GetUserWhereName(ctx context.Context, userName string) (user model.User, err error)
	GetUserWhereNameResetKey(ctx context.Context, userName string, resetKey string, requestedAt time.Time) (user model.User, err error)
	UpdatePasswordWhereName(ctx context.Context, password string, userName string) (err error)
	UpdateIconWhereName(ctx context.Context, iconURL string, userName string) (err error)
	UpdateSelfIntroductionWhereName(ctx context.Context, selfIntroduction string, userName string) (err error)
	UpdateEmailVerifiedWhereNameActivationKey(ctx context.Context, emailVerified bool, userName string, activationKey string) (err error)
	UpdatePasswordResetWhereEmail(ctx context.Context, count uint16, resetKey string, expiresAt time.Time, email string) (err error)
	ResetPassword(ctx context.Context, password string, userName string, resetKey string) (err error)
	UpdatePostingCount(ctx context.Context, userName string, increment bool) (err error)
	UpdateLikeCount(ctx context.Context, userName string, increment bool) (err error)
	UpdateLikedCount(ctx context.Context, userName string) (err error)
	UpdateLikedCountDecrementWhenUserDelete(ctx context.Context, userName string) (err error)
	UpdateFollowCount(ctx context.Context, userName string, increment bool) (err error)
	UpdateFollowedCount(ctx context.Context, userName string, increment bool) (err error)
	UpdateFollowCountDecrementWhereFollowedUserName(ctx context.Context, userName string) (err error)
	UpdateFollowedCountDecrementWhereFollowingUserName(ctx context.Context, userName string) (err error)
	DeleteWhereName(ctx context.Context, userName string) (err error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (err error) {
	q := "INSERT INTO `users` (`name`, `email`, `password`, `activation_key`) VALUES (?, ?, ?, ?)"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, user.Name, user.Email, user.Password, user.ActivationKey)
	} else {
		_, err = r.db.ExecContext(ctx, q, user.Name, user.Email, user.Password, user.ActivationKey)
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok && mysqlErr.Number == 1062 {
		return ErrDuplicateData
	}
	return
}

func (r *UserRepository) GetUserWhereEmail(ctx context.Context, email string) (user model.User, err error) {
	q := "SELECT `name`, `email`, `password`, `icon`, `self_introduction`, `posting_count`, `like_count`, `liked_count`, `follow_count`, `followed_count`, `activation_key`, `email_verified`, `password_reset_email_count`, `password_reset_key`, `password_reset_key_expires_at`, `created_at`, `updated_at` FROM `users` WHERE `email` = ?"
	err = r.db.QueryRowContext(ctx, q, email).Scan(&user.Name, &user.Email, &user.Password, &user.Icon, &user.SelfIntroduction, &user.PostingCount, &user.LikeCount, &user.LikedCount, &user.FollowCount, &user.FollowedCount, &user.ActivationKey, &user.EmailVerified, &user.PasswordResetEmailCount, &user.PasswordResetKey, &user.PasswordResetKeyExpiresAt, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		err = ErrNotExistsData
		return
	}
	return
}

func (r *UserRepository) GetUserWhereName(ctx context.Context, userName string) (user model.User, err error) {
	q := "SELECT `name`, `email`, `password`, `icon`, `self_introduction`, `posting_count`, `like_count`, `liked_count`, `follow_count`, `followed_count`, `activation_key`, `email_verified`, `password_reset_email_count`, `password_reset_key`, `password_reset_key_expires_at`, `created_at`, `updated_at` FROM `users` WHERE `name` = ?"
	err = r.db.QueryRowContext(ctx, q, userName).Scan(&user.Name, &user.Email, &user.Password, &user.Icon, &user.SelfIntroduction, &user.PostingCount, &user.LikeCount, &user.LikedCount, &user.FollowCount, &user.FollowedCount, &user.ActivationKey, &user.EmailVerified, &user.PasswordResetEmailCount, &user.PasswordResetKey, &user.PasswordResetKeyExpiresAt, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		err = ErrNotExistsData
		return
	}
	return
}

func (r *UserRepository) GetUserWhereNameResetKey(ctx context.Context, userName, resetKey string, requestedAt time.Time) (user model.User, err error) {
	q := "SELECT `name`, `email`, `password`, `icon`, `self_introduction`, `posting_count`, `like_count`, `liked_count`, `follow_count`, `followed_count`, `activation_key`, `email_verified`, `password_reset_email_count`, `password_reset_key`, `password_reset_key_expires_at`, `created_at`, `updated_at` FROM `users` WHERE `name` = ? AND `password_reset_key` = ? AND `password_reset_key_expires_at` >= ?"
	err = r.db.QueryRowContext(ctx, q, userName, resetKey, requestedAt).Scan(&user.Name, &user.Email, &user.Password, &user.Icon, &user.SelfIntroduction, &user.PostingCount, &user.LikeCount, &user.LikedCount, &user.FollowCount, &user.FollowedCount, &user.ActivationKey, &user.EmailVerified, &user.PasswordResetEmailCount, &user.PasswordResetKey, &user.PasswordResetKeyExpiresAt, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		err = ErrNotExistsData
		return
	}
	return
}

func (r *UserRepository) UpdatePasswordWhereName(ctx context.Context, password, userName string) (err error) {
	q := "UPDATE `users` SET `password` = ? WHERE `name` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, password, userName)
	} else {
		_, err = r.db.ExecContext(ctx, q, password, userName)
	}
	return
}

func (r *UserRepository) UpdateIconWhereName(ctx context.Context, iconURL, userName string) (err error) {
	q := "UPDATE `users` SET `icon` = ? WHERE `name` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, iconURL, userName)
	} else {
		_, err = r.db.ExecContext(ctx, q, iconURL, userName)
	}
	return
}

func (r *UserRepository) UpdateSelfIntroductionWhereName(ctx context.Context, selfIntroduction, userName string) (err error) {
	q := "UPDATE `users` SET `self_introduction` = ? WHERE `name` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, selfIntroduction, userName)
	} else {
		_, err = r.db.ExecContext(ctx, q, selfIntroduction, userName)
	}
	return
}

func (r *UserRepository) UpdateEmailVerifiedWhereNameActivationKey(ctx context.Context, emailVerified bool, userName, activationKey string) (err error) {
	q := "UPDATE `users` SET `email_verified` = ? WHERE `name` = ? AND `activation_key` = ?"
	tx := m.GetTransaction(ctx)
	var result sql.Result
	if tx != nil {
		result, err = tx.ExecContext(ctx, q, emailVerified, userName, activationKey)
	} else {
		result, err = r.db.ExecContext(ctx, q, emailVerified, userName, activationKey)
	}
	if err != nil {
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return
	}
	if rows == 0 {
		return ErrUserActivationNotFound
	}
	return
}

func (r *UserRepository) UpdatePasswordResetWhereEmail(ctx context.Context, count uint8, resetKey string, expiresAt time.Time, email string) (err error) {
	q := "UPDATE `users` SET `password_reset_email_count` = ?, `password_reset_key` = ?, `password_reset_key_expires_at` = ? WHERE `email` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, count, resetKey, expiresAt, email)
	} else {
		_, err = r.db.ExecContext(ctx, q, count, resetKey, expiresAt, email)
	}
	return
}

func (r *UserRepository) ResetPassword(ctx context.Context, password, userName, resetKey string) (err error) {
	q := "UPDATE `users` SET `password` = ? WHERE `name` = ? AND `password_reset_key` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, password, userName, resetKey)
	} else {
		_, err = r.db.ExecContext(ctx, q, password, userName, resetKey)
	}
	return
}

func (r *UserRepository) UpdateLikeCount(ctx context.Context, userName string, increment bool) (err error) {
	var q string
	if increment {
		q = "UPDATE `users` SET `like_count` = `like_count` + 1 WHERE `name` = ?"
	} else {
		q = "UPDATE `users` SET `like_count` = `like_count` - 1 WHERE `name` = ?"
	}
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userName)
	} else {
		_, err = r.db.ExecContext(ctx, q, userName)
	}
	return
}

func (r *UserRepository) UpdatePostingCount(ctx context.Context, userName string, increment bool) (err error) {
	var q string
	if increment {
		q = "UPDATE `users` SET `posting_count` = `posting_count` + 1 WHERE `name` = ?"
	} else {
		q = "UPDATE `users` SET `posting_count` = `posting_count` - 1 WHERE `name` = ?"
	}
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userName)
	} else {
		_, err = r.db.ExecContext(ctx, q, userName)
	}
	return
}

func (r *UserRepository) UpdateLikedCount(ctx context.Context, postingID int64, increment bool) (err error) {
	var q string
	if increment {
		q = "UPDATE `users` SET `liked_count` = `liked_count` + 1 WHERE `name` = (SELECT `user_name` FROM `postings` WHERE `id` = ?)"
	} else {
		q = "UPDATE `users` SET `liked_count` = `liked_count` - 1 WHERE `name` = (SELECT `user_name` FROM `postings` WHERE `id` = ?)"
	}
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, postingID)
	} else {
		_, err = r.db.ExecContext(ctx, q, postingID)
	}
	return
}

func (r *UserRepository) UpdateLikedCountDecrementWhenUserDelete(ctx context.Context, userName string) (err error) {
	q := "UPDATE `users` SET `liked_count` = `liked_count` - 1 WHERE `name` IN (SELECT `user_name` FROM `postings` WHERE `id` IN (SELECT `posting_id` FROM `likes` WHERE `user_name` = ?))"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userName)
	} else {
		_, err = r.db.ExecContext(ctx, q, userName)
	}
	return
}

func (r *UserRepository) UpdateFollowCount(ctx context.Context, userName string, increment bool) (err error) {
	var q string
	if increment {
		q = "UPDATE `users` SET `follow_count` = `follow_count` + 1 WHERE `name` = ?"
	} else {
		q = "UPDATE `users` SET `follow_count` = `follow_count` - 1 WHERE `name` = ?"
	}
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userName)
	} else {
		_, err = r.db.ExecContext(ctx, q, userName)
	}
	return
}

func (r *UserRepository) UpdateFollowedCount(ctx context.Context, userName string, increment bool) (err error) {
	var q string
	if increment {
		q = "UPDATE `users` SET `followed_count` = `followed_count` + 1 WHERE `name` = ?"
	} else {
		q = "UPDATE `users` SET `followed_count` = `followed_count` - 1 WHERE `name` = ?"
	}
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userName)
	} else {
		_, err = r.db.ExecContext(ctx, q, userName)
	}
	return
}

func (r *UserRepository) UpdateFollowCountDecrementWhereFollowedUserName(ctx context.Context, userName string) (err error) {
	q := "UPDATE `users` SET `follow_count` = `follow_count` - 1 WHERE `name` IN (SELECT `following_user_name` FROM `follows` WHERE `followed_user_name` = ?)"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userName)
	} else {
		_, err = r.db.ExecContext(ctx, q, userName)
	}
	return
}

func (r *UserRepository) UpdateFollowedCountDecrementWhereFollowingUserName(ctx context.Context, userName string) (err error) {
	q := "UPDATE `users` SET `followed_count` = `followed_count` - 1 WHERE `name` IN (SELECT `followed_user_name` FROM `follows` WHERE `following_user_name` = ?)"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userName)
	} else {
		_, err = r.db.ExecContext(ctx, q, userName)
	}
	return
}

func (r *UserRepository) DeleteWhereName(ctx context.Context, userName string) (err error) {
	q := "DELETE FROM `users` WHERE `name` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userName)
	} else {
		_, err = r.db.ExecContext(ctx, q, userName)
	}
	return
}
