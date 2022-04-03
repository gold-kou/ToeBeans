package repository

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"

	m "github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
)

type NotificationRepositoryInterface interface {
	Create(ctx context.Context, notification *model.Notification) (err error)
	GetNotifications(ctx context.Context, userID int64) (notifications []model.Notification, err error)
	DeleteWhereID(ctx context.Context, id int64) (err error)
	DeleteWhereNotificationUserID(ctx context.Context, userID int64) (err error)
}

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) Create(ctx context.Context, notification *model.Notification) (err error) {
	q := "INSERT INTO `notifications` (`visitor_user_id`, `visited_user_id`, `action`) VALUES (?, ?, ?, ?, ?)"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, notification.VisitorUserID, notification.VisitedUserID, notification.Action)
	} else {
		_, err = r.db.ExecContext(ctx, q, notification.VisitorUserID, notification.VisitedUserID, notification.Action)
	}
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok && mysqlErr.Number == 1062 {
		return ErrDuplicateData
	}
	return
}

func (r *NotificationRepository) GetNotifications(ctx context.Context, userID int64) (notifications []model.Notification, err error) {
	q := "SELECT `id`, `visitor_user_id`, `visited_user_id`, `action`, `created_at`, `updated_at` FROM `postings` WHERE `created_at` < ? ORDER BY `created_at` DESC"
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err == sql.ErrNoRows {
		err = ErrNotExistsData
		return
	}
	if err != nil {
		return
	}
	defer rows.Close()

	var n model.Notification
	for rows.Next() {
		if err = rows.Scan(&n.ID, &n.VisitorUserID, &n.VisitedUserID, &n.Action, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return
		}
		notifications = append(notifications, n)
		n = model.Notification{}
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (r *NotificationRepository) DeleteWhereID(ctx context.Context, id int64) (err error) {
	q := "DELETE FROM `notifications` WHERE `id` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, id)
	} else {
		_, err = r.db.ExecContext(ctx, q, id)
	}
	return
}

func (r *NotificationRepository) DeleteWhereVisitedName(ctx context.Context, userID int64) (err error) {
	q := "DELETE FROM `notifications` WHERE `visited_user_id` = ?"
	tx := m.GetTransaction(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, q, userID)
	} else {
		_, err = r.db.ExecContext(ctx, q, userID)
	}
	return
}
