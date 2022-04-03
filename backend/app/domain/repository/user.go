package repository

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"

	m "github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *model.User) (err error)
	GetUserWhereID(ctx context.Context, id int64) (user model.User, err error)
	GetUserWhereName(ctx context.Context, userName string) (user model.User, err error)
	GetUserWhereEmail(ctx context.Context, email string) (user model.User, err error)
	UpdatePasswordWhereName(ctx context.Context, password string, userName string) (err error)
	UpdateIconWhereName(ctx context.Context, iconURL string, userName string) (err error)
	UpdateSelfIntroductionWhereName(ctx context.Context, selfIntroduction string, userName string) (err error)
	UpdateEmailVerifiedWhereNameActivationKey(ctx context.Context, emailVerified bool, userName string, activationKey string) (err error)
	ResetPassword(ctx context.Context, password string, userName string) (err error)
	DeleteWhereID(ctx context.Context, id int64) (err error)
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

func (r *UserRepository) GetUserWhereID(ctx context.Context, id int64) (user model.User, err error) {
	q := "SELECT `id`, `name`, `email`, `password`, `icon`, `self_introduction`, `activation_key`, `email_verified`, `created_at`, `updated_at` FROM `users` WHERE `id` = ?"
	err = r.db.QueryRowContext(ctx, q, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Icon, &user.SelfIntroduction, &user.ActivationKey, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		err = ErrNotExistsData
		return
	}
	return
}

func (r *UserRepository) GetUserWhereName(ctx context.Context, userName string) (user model.User, err error) {
	q := "SELECT `id`, `name`, `email`, `password`, `icon`, `self_introduction`, `activation_key`, `email_verified`, `created_at`, `updated_at` FROM `users` WHERE `name` = ?"
	err = r.db.QueryRowContext(ctx, q, userName).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Icon, &user.SelfIntroduction, &user.ActivationKey, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		err = ErrNotExistsData
		return
	}
	return
}

func (r *UserRepository) GetUserWhereEmail(ctx context.Context, email string) (user model.User, err error) {
	q := "SELECT `id`, `name`, `email`, `password`, `icon`, `self_introduction`, `activation_key`, `email_verified`, `created_at`, `updated_at` FROM `users` WHERE `email` = ?"
	err = r.db.QueryRowContext(ctx, q, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Icon, &user.SelfIntroduction, &user.ActivationKey, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
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

func (r *UserRepository) ResetPassword(ctx context.Context, password, userName string) (err error) {
	q := "UPDATE `users` SET `password` = ? WHERE `name` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, password, userName)
	} else {
		_, err = r.db.ExecContext(ctx, q, password, userName)
	}
	return
}

func (r *UserRepository) DeleteWhereID(ctx context.Context, id int64) (err error) {
	q := "DELETE FROM `users` WHERE `id` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, id)
	} else {
		_, err = r.db.ExecContext(ctx, q, id)
	}
	return
}
