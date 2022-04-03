package repository

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"

	m "github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

type CommentRepositoryInterface interface {
	Create(ctx context.Context, comment *model.Comment) (err error)
	GetCommentsWherePostingID(ctx context.Context, id int64) (comments []model.Comment, err error)
	GetWhereID(ctx context.Context, id int64) (comment model.Comment, err error)
	DeleteWhereID(ctx context.Context, id int64) (err error)
	DeleteWhereUserID(ctx context.Context, userID int64) (err error)
}

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) Create(ctx context.Context, comment *model.Comment) (err error) {
	q := "INSERT INTO `comments` (`user_id`, `posting_id`, `comment`) VALUES (?, ?, ?)"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, comment.UserID, comment.PostingID, comment.Comment)
	} else {
		_, err = r.db.ExecContext(ctx, q, comment.UserID, comment.PostingID, comment.Comment)
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok && mysqlErr.Number == 1062 {
		return ErrDuplicateData
	}
	return
}

func (r *CommentRepository) GetCommentsWherePostingID(ctx context.Context, postingID int64) (comments []model.Comment, err error) {
	q := "SELECT `id`, `user_id`, `posting_id`, `comment`, `created_at`, `updated_at` FROM `comments` WHERE `posting_id` = ? ORDER BY `created_at` DESC"
	rows, err := r.db.QueryContext(ctx, q, postingID)
	if err == sql.ErrNoRows {
		err = ErrNotExistsData
		return
	}
	if err != nil {
		return
	}
	defer rows.Close()

	var c model.Comment
	for rows.Next() {
		if err = rows.Scan(&c.ID, &c.UserID, &c.PostingID, &c.Comment, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return
		}
		comments = append(comments, c)
		c = model.Comment{}
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (r *CommentRepository) GetWhereID(ctx context.Context, id int64) (comment model.Comment, err error) {
	q := "SELECT `id`, `user_id`, `posting_id`, `created_at`, `updated_at` FROM `comments` WHERE `id` = ?"
	err = r.db.QueryRowContext(ctx, q, id).Scan(&comment.ID, &comment.UserID, &comment.PostingID, &comment.Comment, &comment.CreatedAt, &comment.UpdatedAt)
	if err == sql.ErrNoRows {
		err = ErrNotExistsData
		return
	}
	return
}

func (r *CommentRepository) DeleteWhereID(ctx context.Context, id int64) (err error) {
	q := "DELETE FROM `comments` WHERE `id` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, id)
	} else {
		_, err = r.db.ExecContext(ctx, q, id)
	}
	return
}

func (r *CommentRepository) DeleteWhereUserID(ctx context.Context, userID int64) (err error) {
	q := "DELETE FROM `comments` WHERE `user_id` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userID)
	} else {
		_, err = r.db.ExecContext(ctx, q, userID)
	}
	return
}
