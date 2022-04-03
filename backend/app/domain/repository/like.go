package repository

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"

	m "github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

type LikeRepositoryInterface interface {
	Create(ctx context.Context, like *model.Like) (err error)
	GetWhereUserID(ctx context.Context, userID int64) (like model.Like, err error)
	GetWhereUserIDPostingID(ctx context.Context, userID int64, postingID int64) (like model.Like, err error)
	GetLikeCountWhereUserID(ctx context.Context, userID int64) (int64, err error)
	GetLikedCountWhereUserID(ctx context.Context, userID int64) (int64, err error)
	GetLikedCountWherePostingID(ctx context.Context, postingID int64) (int64, err error)
	DeleteWhereUserIDPostingID(ctx context.Context, userID int64, postingID int64) (err error)
	DeleteWhereUserID(ctx context.Context, userID int64) (err error)
	DeleteWhereInPosingIDs(ctx context.Context, userID int64) (err error)
}

type LikeRepository struct {
	db *sql.DB
}

func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{
		db: db,
	}
}

func (r *LikeRepository) Create(ctx context.Context, like *model.Like) (err error) {
	q := "INSERT INTO `likes` (`user_id`, `posting_id`) VALUES (?, ?)"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, like.UserID, like.PostingID)
	} else {
		_, err = r.db.ExecContext(ctx, q, like.UserID, like.PostingID)
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok && mysqlErr.Number == 1062 {
		return ErrDuplicateData
	}
	return
}

func (r *LikeRepository) GetWhereUserID(ctx context.Context, userID int64) (likes []model.Like, err error) {
	q := "SELECT `id`, `user_id`, `posting_id`, `created_at`, `updated_at` FROM `likes` WHERE `user_id` = ?"
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err == sql.ErrNoRows {
		err = ErrNotExistsData
		return
	}
	if err != nil {
		return
	}
	defer rows.Close()

	var like model.Like
	for rows.Next() {
		if err = rows.Scan(&like.ID, &like.UserID, &like.PostingID, &like.CreatedAt, &like.UpdatedAt); err != nil {
			return
		}
		likes = append(likes, like)
		like = model.Like{}
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (r *LikeRepository) GetWhereUserIDPostingID(ctx context.Context, userID, postingID int64) (like model.Like, err error) {
	q := "SELECT `id`, `user_id`, `posting_id`, `created_at`, `updated_at` FROM `likes` WHERE `user_id` = ? AND `posting_id` = ?"
	err = r.db.QueryRowContext(ctx, q, userID, postingID).Scan(&like.ID, &like.UserID, &like.PostingID, &like.CreatedAt, &like.UpdatedAt)
	if err == sql.ErrNoRows {
		err = ErrNotExistsData
		return
	}
	return
}

func (r *LikeRepository) GetLikeCountWhereUserID(ctx context.Context, userID int64) (count int64, err error) {
	q := "SELECT COUNT(*) FROM `likes` WHERE `user_id` = ?"
	err = r.db.QueryRowContext(ctx, q, userID).Scan(&count)
	return
}

func (r *LikeRepository) GetLikedCountWhereUserID(ctx context.Context, userID int64) (count int64, err error) {
	q := "SELECT COUNT(*) FROM `likes` WHERE `posting_id` IN(SELECT `id` FROM `postings` WHERE `user_id` = ?);"
	err = r.db.QueryRowContext(ctx, q, userID).Scan(&count)
	return
}

func (r *LikeRepository) GetLikedCountWherePostingID(ctx context.Context, postingID int64) (count int64, err error) {
	q := "SELECT COUNT(*) FROM `likes` WHERE `posting_id` = ?;"
	err = r.db.QueryRowContext(ctx, q, postingID).Scan(&count)
	return
}

func (r *LikeRepository) DeleteWhereUserIDPostingID(ctx context.Context, userID, postingID int64) (err error) {
	q := "DELETE FROM `likes` WHERE `user_id` = ? AND `posting_id` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userID, postingID)
	} else {
		_, err = r.db.ExecContext(ctx, q, userID, postingID)
	}
	return
}

func (r *LikeRepository) DeleteWhereUserID(ctx context.Context, userID int64) (err error) {
	q := "DELETE FROM `likes` WHERE `user_id` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userID)
	} else {
		_, err = r.db.ExecContext(ctx, q, userID)
	}
	return
}

func (r *LikeRepository) DeleteWhereInPosingIDs(ctx context.Context, userID int64) (err error) {
	q := "DELETE FROM `likes` WHERE `posting_id` IN (SELECT `id` FROM `postings` WHERE `user_id` = ?)"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userID)
	} else {
		_, err = r.db.ExecContext(ctx, q, userID)
	}
	return
}
