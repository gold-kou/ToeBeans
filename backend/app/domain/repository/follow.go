package repository

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"

	m "github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

type FollowRepositoryInterface interface {
	FindByBothUserIDs(ctx context.Context, followingUserID, followedUserID int64) (follow model.Follow, err error)
	GetFollowCountWhereUserID(ctx context.Context, userID int64) (int64, err error)
	GetFollowedCountWhereUserID(ctx context.Context, userID int64) (int64, err error)
	Create(ctx context.Context, follow *model.Follow) (err error)
	DeleteWhereBothUserIDs(ctx context.Context, followingUserID, followedUserID int64) (err error)
	DeleteWhereFollowingUserID(ctx context.Context, userID int64) (err error)
	DeleteWhereFollowedUserID(ctx context.Context, userID int64) (err error)
}

type FollowRepository struct {
	db *sql.DB
}

func NewFollowRepository(db *sql.DB) *FollowRepository {
	return &FollowRepository{
		db: db,
	}
}

func (r *FollowRepository) GetFollowCountWhereUserID(ctx context.Context, userID int64) (count int64, err error) {
	q := "SELECT COUNT(*) FROM `follows` WHERE `following_user_id` = ?"
	err = r.db.QueryRowContext(ctx, q, userID).Scan(&count)
	return
}

func (r *FollowRepository) GetFollowedCountWhereUserID(ctx context.Context, userID int64) (count int64, err error) {
	q := "SELECT COUNT(*) FROM `follows` WHERE `followed_user_id` = ?"
	err = r.db.QueryRowContext(ctx, q, userID).Scan(&count)
	return
}

func (r *FollowRepository) Create(ctx context.Context, follow *model.Follow) (err error) {
	q := "INSERT INTO `follows` (`following_user_id`, `followed_user_id`) VALUES (?, ?)"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, follow.FollowingUserID, follow.FollowedUserID)
	} else {
		_, err = r.db.ExecContext(ctx, q, follow.FollowingUserID, follow.FollowedUserID)
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok && mysqlErr.Number == 1062 {
		return ErrDuplicateData
	}
	return
}

func (r *FollowRepository) FindByBothUserIDs(ctx context.Context, followingUserID, followedUserID int64) (follow model.Follow, err error) {
	q := "SELECT `id`, `following_user_id`, `followed_user_id`, `created_at`, `updated_at` FROM `follows` WHERE `following_user_id` = ? AND `followed_user_id` = ?"
	err = r.db.QueryRowContext(ctx, q, followingUserID, followedUserID).Scan(&follow.ID, &follow.FollowingUserID, &follow.FollowedUserID, &follow.CreatedAt, &follow.UpdatedAt)
	if err == sql.ErrNoRows {
		err = ErrNotExistsData
		return
	}
	return
}

func (r *FollowRepository) DeleteWhereBothUserIDs(ctx context.Context, followingUserID, followedUserID int64) (err error) {
	q := "DELETE FROM `follows` WHERE `following_user_id` = ? AND `followed_user_id` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, followingUserID, followedUserID)
	} else {
		_, err = r.db.ExecContext(ctx, q, followingUserID, followedUserID)
	}
	return
}

func (r *FollowRepository) DeleteWhereFollowingUserID(ctx context.Context, userID int64) (err error) {
	q := "DELETE FROM `follows` WHERE `following_user_id` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userID)
	} else {
		_, err = r.db.ExecContext(ctx, q, userID)
	}
	return
}

func (r *FollowRepository) DeleteWhereFollowedUserID(ctx context.Context, userID int64) (err error) {
	q := "DELETE FROM `follows` WHERE `followed_user_id` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userID)
	} else {
		_, err = r.db.ExecContext(ctx, q, userID)
	}
	return
}
