package repository

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"

	m "github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

type FollowRepositoryInterface interface {
	Create(ctx context.Context, follow *model.Follow) (err error)
	DeleteWhereFollowingFollowedUserName(ctx context.Context, followingUserName, followedUserName string) (err error)
	DeleteWhereFollowingUserName(ctx context.Context, userName string) (err error)
}

type FollowRepository struct {
	db *sql.DB
}

func NewFollowRepository(db *sql.DB) *FollowRepository {
	return &FollowRepository{
		db: db,
	}
}

func (r *FollowRepository) Create(ctx context.Context, follow *model.Follow) (err error) {
	q := "INSERT INTO `follows` (`following_user_name`, `followed_user_name`) VALUES (?, ?)"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, follow.FollowingUserName, follow.FollowedUserName)
	} else {
		_, err = r.db.ExecContext(ctx, q, follow.FollowingUserName, follow.FollowedUserName)
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok && mysqlErr.Number == 1062 {
		return ErrDuplicateData
	}
	return
}

func (r *FollowRepository) DeleteWhereFollowingFollowedUserName(ctx context.Context, followingUserName, followedUserName string) (err error) {
	q := "DELETE FROM `follows` WHERE `following_user_name` = ? AND `followed_user_name` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, followingUserName, followedUserName)
	} else {
		_, err = r.db.ExecContext(ctx, q, followingUserName, followedUserName)
	}
	return
}

func (r *FollowRepository) DeleteWhereFollowingUserName(ctx context.Context, userName string) (err error) {
	q := "DELETE FROM `follows` WHERE `following_user_name` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userName)
	} else {
		_, err = r.db.ExecContext(ctx, q, userName)
	}
	return
}
