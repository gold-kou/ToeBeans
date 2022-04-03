package repository

import (
	"context"
	"database/sql"
	"time"

	m "github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

type PostingRepositoryInterface interface {
	Create(ctx context.Context, posting *model.Posting) (err error)
	GetPostings(ctx context.Context, sinceAt time.Time, limit int8) (postings []model.Posting, err error)
	GetWhereID(ctx context.Context, id int64) (posting model.Posting, err error)
	GetWhereIDUserID(ctx context.Context, id int64, userID int64) (posting model.Posting, err error)
	GetCountWhereUserID(ctx context.Context, userID int64) (int64, err error)
	DeleteWhereID(ctx context.Context, id int64) (err error)
	DeleteWhereUserID(ctx context.Context, userID int64) (err error)
}

type PostingRepository struct {
	db *sql.DB
}

func NewPostingRepository(db *sql.DB) *PostingRepository {
	return &PostingRepository{
		db: db,
	}
}

func (r *PostingRepository) Create(ctx context.Context, posting *model.Posting) (err error) {
	q := "INSERT INTO `postings` (`user_id`, `title`, `image_url`) VALUES (?, ?, ?)"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, posting.UserID, posting.Title, posting.ImageURL)
	} else {
		_, err = r.db.ExecContext(ctx, q, posting.UserID, posting.Title, posting.ImageURL)
	}
	return
}

func (r *PostingRepository) GetPostings(ctx context.Context, sinceAt time.Time, limit int8, userID int64) (postings []model.Posting, err error) {
	var q string
	var rows *sql.Rows
	if userID == 0 {
		q = "SELECT `id`, `user_id`, `title`, `image_url`, `created_at`, `updated_at` FROM `postings` WHERE `created_at` < ? ORDER BY `created_at` DESC LIMIT ?"
		rows, err = r.db.QueryContext(ctx, q, sinceAt, limit)
	} else {
		q = "SELECT `id`, `user_id`, `title`, `image_url`, `created_at`, `updated_at` FROM `postings` WHERE `created_at` < ? AND `user_id` = ? ORDER BY `created_at` DESC LIMIT ?"
		rows, err = r.db.QueryContext(ctx, q, sinceAt, userID, limit)
	}
	if err == sql.ErrNoRows {
		err = ErrNotExistsData
		return
	}
	if err != nil {
		return
	}
	defer rows.Close()

	var p model.Posting
	for rows.Next() {
		if err = rows.Scan(&p.ID, &p.UserID, &p.Title, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return
		}
		postings = append(postings, p)
		p = model.Posting{}
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (r *PostingRepository) GetWhereID(ctx context.Context, id int64) (posting model.Posting, err error) {
	q := "SELECT `id`, `user_id`, `title`, `image_url`, `created_at`, `updated_at` FROM `postings` WHERE `id` = ?"
	err = r.db.QueryRowContext(ctx, q, id).Scan(&posting.ID, &posting.UserID, &posting.Title, &posting.ImageURL, &posting.CreatedAt, &posting.UpdatedAt)
	if err == sql.ErrNoRows {
		err = ErrNotExistsData
		return
	}
	return
}

func (r *PostingRepository) GetWhereIDUserID(ctx context.Context, id int64, userID int64) (posting model.Posting, err error) {
	q := "SELECT `id`, `user_id`, `title`, `image_url`, `created_at`, `updated_at` FROM `postings` WHERE `id` = ? AND `user_id` = ?"
	err = r.db.QueryRowContext(ctx, q, id, userID).Scan(&posting.ID, &posting.UserID, &posting.Title, &posting.ImageURL, &posting.CreatedAt, &posting.UpdatedAt)
	if err == sql.ErrNoRows {
		err = ErrNotExistsData
		return
	}
	return
}

func (r *PostingRepository) GetCountWhereUserID(ctx context.Context, userID int64) (count int64, err error) {
	q := "SELECT COUNT(*) FROM `postings` WHERE `user_id` = ?"
	err = r.db.QueryRowContext(ctx, q, userID).Scan(&count)
	return
}

func (r *PostingRepository) DeleteWhereID(ctx context.Context, id int64) (err error) {
	q := "DELETE FROM `postings` WHERE `id` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, id)
	} else {
		_, err = r.db.ExecContext(ctx, q, id)
	}
	return
}

func (r *PostingRepository) DeleteWhereUserID(ctx context.Context, userID int64) (err error) {
	q := "DELETE FROM `postings` WHERE `user_id` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userID)
	} else {
		_, err = r.db.ExecContext(ctx, q, userID)
	}
	return
}
