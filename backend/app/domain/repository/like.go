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
	GetWhereUserName(ctx context.Context, userName string) (like model.Like, err error)
	GetWhereID(ctx context.Context, id int64) (like model.Like, err error)
	DeleteWhereID(ctx context.Context, id int64) (err error)
	DeleteWhereUserName(ctx context.Context, userName string) (err error)
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
	q := "INSERT INTO `likes` (`user_name`, `posting_id`) VALUES (?, ?)"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, like.UserName, like.PostingID)
	} else {
		_, err = r.db.ExecContext(ctx, q, like.UserName, like.PostingID)
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok && mysqlErr.Number == 1062 {
		return ErrDuplicateData
	}
	return
}

func (r *LikeRepository) GetWhereUserName(ctx context.Context, userName string) (likes []model.Like, err error) {
	q := "SELECT `id`, `user_name`, `posting_id`, `created_at`, `updated_at` FROM `likes` WHERE `user_name` = ?"
	rows, err := r.db.QueryContext(ctx, q, userName)
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
		if err = rows.Scan(&like.ID, &like.UserName, &like.PostingID, &like.CreatedAt, &like.UpdatedAt); err != nil {
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

func (r *LikeRepository) GetWhereID(ctx context.Context, id int64) (like model.Like, err error) {
	q := "SELECT `id`, `user_name`, `posting_id`, `created_at`, `updated_at` FROM `likes` WHERE `id` = ?"
	err = r.db.QueryRowContext(ctx, q, id).Scan(&like.ID, &like.UserName, &like.PostingID, &like.CreatedAt, &like.UpdatedAt)
	if err == sql.ErrNoRows {
		err = ErrNotExistsData
		return
	}
	return
}

func (r *LikeRepository) DeleteWhereID(ctx context.Context, id int64) (err error) {
	q := "DELETE FROM `likes` WHERE `id` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, id)
	} else {
		_, err = r.db.ExecContext(ctx, q, id)
	}
	return
}

func (r *LikeRepository) DeleteWhereUserName(ctx context.Context, userName string) (err error) {
	q := "DELETE FROM `likes` WHERE `user_name` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userName)
	} else {
		_, err = r.db.ExecContext(ctx, q, userName)
	}
	return
}
