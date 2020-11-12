package repository

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"

	m "github.com/gold-kou/ToeBeans/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/app/domain/model"
)

type CommentRepositoryInterface interface {
	Create(ctx context.Context, comment *model.Comment) (err error)
	GetCommentsWherePostingID(ctx context.Context, id int64) (comments []model.Comment, err error)
	GetWhereID(ctx context.Context, id int64) (comment model.Comment, err error)
	DeleteWhereID(ctx context.Context, id int64) (err error)
	DeleteWhereUserName(ctx context.Context, userName string) (err error)
}

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) Create(ctx context.Context, comment *model.Comment) (result sql.Result, err error) {
	q := "INSERT INTO `comments` (`user_name`, `posting_id`, `comment`) VALUES (?, ?, ?)"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		result, err = tx.ExecContext(ctx, q, comment.UserName, comment.PostingID, comment.Comment)
	} else {
		result, err = r.db.ExecContext(ctx, q, comment.UserName, comment.PostingID, comment.Comment)
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok && mysqlErr.Number == 1062 {
		return nil, ErrDuplicateData
	}
	return
}

func (r *CommentRepository) GetCommentsWherePostingID(ctx context.Context, postingID int64) (comments []model.Comment, err error) {
	q := "SELECT `id`, `user_name`, `posting_id`, `comment`, `created_at`, `updated_at` FROM `comments` WHERE `posting_id` = ? ORDER BY `created_at` DESC"
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
		if err = rows.Scan(&c.ID, &c.UserName, &c.PostingID, &c.Comment, &c.CreatedAt, &c.UpdatedAt); err != nil {
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
	q := "SELECT `id`, `user_name`, `posting_id`, `created_at`, `updated_at` FROM `comments` WHERE `id` = ?"
	err = r.db.QueryRowContext(ctx, q, id).Scan(&comment.ID, &comment.UserName, &comment.PostingID, &comment.Comment, &comment.CreatedAt, &comment.UpdatedAt)
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

func (r *CommentRepository) DeleteWhereUserName(ctx context.Context, userName string) (err error) {
	q := "DELETE FROM `comments` WHERE `user_name` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userName)
	} else {
		_, err = r.db.ExecContext(ctx, q, userName)
	}
	return
}
