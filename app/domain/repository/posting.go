package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"

	m "github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/model"
)

type PostingRepositoryInterface interface {
	Create(ctx context.Context, posting *model.Posting) (err error)
	GetPostings(ctx context.Context, sinceAt time.Time, limit int8) (postings []model.Posting, err error)
	GetWhereIDUserName(ctx context.Context, id uint64, userName string) (posting model.Posting, err error)
	UpdateLikedCount(ctx context.Context, id int64, increment bool) (err error)
	DeleteWhereID(ctx context.Context, id uint64) (err error)
	DeleteWhereUserName(ctx context.Context, userName string) (err error)
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
	q := "INSERT INTO `postings` (`user_name`, `title`, `image_url`) VALUES (?, ?, ?)"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, posting.UserName, posting.Title, posting.ImageURL)
	} else {
		_, err = r.db.ExecContext(ctx, q, posting.UserName, posting.Title, posting.ImageURL)
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok && mysqlErr.Number == 1062 {
		return ErrDuplicateData
	}
	return
}

func (r *PostingRepository) GetPostings(ctx context.Context, sinceAt time.Time, limit int8, userName string) (postings []model.Posting, err error) {
	var q string
	var rows *sql.Rows
	if userName == "" {
		q = "SELECT `id`, `user_name`, `title`, `image_url`, `liked_count`, `created_at` FROM `postings` WHERE `created_at` < ? ORDER BY `created_at` DESC LIMIT ?"
		rows, err = r.db.QueryContext(ctx, q, sinceAt, limit)
	} else {
		q = "SELECT `id`, `user_name`, `title`, `image_url`, `liked_count`, `created_at` FROM `postings` WHERE `created_at` < ? AND `user_name` = ? ORDER BY `created_at` DESC LIMIT ?"
		rows, err = r.db.QueryContext(ctx, q, sinceAt, userName, limit)
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
		if err = rows.Scan(&p.ID, &p.UserName, &p.Title, &p.ImageURL, &p.LikedCount, &p.CreatedAt); err != nil {
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

func (r *PostingRepository) GetWhereIDUserName(ctx context.Context, id uint64, userName string) (posting model.Posting, err error) {
	q := "SELECT `id`, `user_name`, `title`, `image_url`, `liked_count`, `created_at`, `updated_at` FROM `postings` WHERE `id` = ? AND `user_name` = ?"
	err = r.db.QueryRowContext(ctx, q, id, userName).Scan(&posting.ID, &posting.UserName, &posting.Title, &posting.ImageURL, &posting.LikedCount, &posting.CreatedAt, &posting.UpdatedAt)
	if err == sql.ErrNoRows {
		err = ErrNotExistsData
		return
	}
	return
}

func (r *PostingRepository) UpdateLikedCount(ctx context.Context, id int64, increment bool) (err error) {
	var q string
	if increment {
		q = "UPDATE `postings` SET `liked_count` = `liked_count` + 1 WHERE `id` = ?"
	} else {
		q = "UPDATE `postings` SET `liked_count` = `liked_count` - 1 WHERE `id` = ?"
	}
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, id)
	} else {
		_, err = r.db.ExecContext(ctx, q, id)
	}
	return
}

func (r *PostingRepository) DeleteWhereIDUserName(ctx context.Context, id uint64, userName string) (err error) {
	q := "DELETE FROM `postings` WHERE `id` = ? AND `user_name` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, id, userName)
	} else {
		_, err = r.db.ExecContext(ctx, q, id, userName)
	}
	return
}

func (r *PostingRepository) DeleteWhereUserName(ctx context.Context, userName string) (err error) {
	q := "DELETE FROM `postings` WHERE `user_name` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userName)
	} else {
		_, err = r.db.ExecContext(ctx, q, userName)
	}
	return
}
