package testing

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/mysql"
	"github.com/gold-kou/ToeBeans/backend/app/domain/model"
	"github.com/gold-kou/ToeBeans/backend/app/lib"
	"github.com/gold-kou/ToeBeans/backend/testing/dummy"
)

func SetTestEnv(key, val string) func() {
	preVal := os.Getenv(key)
	os.Setenv(key, val)
	return func() {
		os.Setenv(key, preVal)
	}
}

func SetTime(t time.Time) {
	lib.NowFunc = func() time.Time { return t }
}

func SetTestTime() {
	SetTime(GetTestTime())
}

func GetTestTime() time.Time {
	loc, _ := time.LoadLocation(os.Getenv("TZ"))
	return time.Date(2020, time.January, 1, 0, 0, 0, 0, loc)
}

func ResetTime() {
	lib.NowFunc = time.Now
}

func SetupDBTest() *sql.DB {
	db, err := mysql.NewDB()
	if err != nil {
		panic(err)
	}
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		db.Close()
		panic(err)
	}
	defer rows.Close()
	tables := []string{}
	for rows.Next() {
		var table string
		if err = rows.Scan(&table); err != nil {
			db.Close()
			panic(err)
		}
		tables = append(tables, table)
	}
	if err = rows.Err(); err != nil {
		db.Close()
		panic(err)
	}

	var tx *sql.Tx
	tx, err = db.Begin()
	if err != nil {
		db.Close()
		panic(err)
	}
	_, err = tx.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		db.Close()
		panic(err)
	}
	for _, table := range tables {
		_, err = tx.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", table))
		if err != nil {
			db.Close()
			panic(err)
		}
	}
	_, err = tx.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		db.Close()
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		db.Close()
		panic(err)
	}

	return db
}

func TeardownDBTest(db *sql.DB) {
	if err := DeleteAllTableData(db, "user_reports"); err != nil {
		panic(err)
	}
	if err := DeleteAllTableData(db, "posting_reports"); err != nil {
		panic(err)
	}
	if err := DeleteAllTableData(db, "notifications"); err != nil {
		panic(err)
	}
	if err := DeleteAllTableData(db, "likes"); err != nil {
		panic(err)
	}
	if err := DeleteAllTableData(db, "comments"); err != nil {
		panic(err)
	}
	if err := DeleteAllTableData(db, "follows"); err != nil {
		panic(err)
	}
	if err := DeleteAllTableData(db, "postings"); err != nil {
		panic(err)
	}
	if err := DeleteAllTableData(db, "password_resets"); err != nil {
		panic(err)
	}
	if err := DeleteAllTableData(db, "users"); err != nil {
		panic(err)
	}
	db.Close()
}

func FindAllUsers(ctx context.Context, db *sql.DB) ([]model.User, error) {
	q := "SELECT `id`, `name`, `email`, `password`, `icon`, `self_introduction`, `activation_key`, `email_verified`, `created_at`, `updated_at` FROM `users`"
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []model.User{}
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Icon, &u.SelfIntroduction, &u.ActivationKey, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, u)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func FindAllPasswordResets(ctx context.Context, db *sql.DB) ([]model.PasswordReset, error) {
	q := "SELECT `id`, `user_id`, `password_reset_email_count`, `password_reset_key`, `password_reset_key_expires_at`, `created_at`, `updated_at` FROM `password_resets`"
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []model.PasswordReset{}
	for rows.Next() {
		var pr model.PasswordReset
		if err := rows.Scan(&pr.ID, &pr.UserID, &pr.PasswordResetEmailCount, &pr.PasswordResetKey, &pr.PasswordResetKeyExpiresAt, &pr.CreatedAt, &pr.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, pr)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func FindAllPostings(ctx context.Context, db *sql.DB) ([]model.Posting, error) {
	q := "SELECT `id`, `user_id`, `title`, `image_url`, `created_at`, `updated_at` FROM `postings`"
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []model.Posting{}
	for rows.Next() {
		var p model.Posting
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func FindAllLikes(ctx context.Context, db *sql.DB) ([]model.Like, error) {
	q := "SELECT `id`, `user_id`, `posting_id`, `created_at`, `updated_at` FROM `likes`"
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []model.Like{}
	for rows.Next() {
		var l model.Like
		if err := rows.Scan(&l.ID, &l.UserID, &l.PostingID, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, l)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func FindAllComments(ctx context.Context, db *sql.DB) ([]model.Comment, error) {
	q := "SELECT `id`, `user_id`, `posting_id`, `comment`, `created_at`, `updated_at` FROM `comments`"
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []model.Comment{}
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.UserID, &c.PostingID, &c.Comment, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, c)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func FindAllFollows(ctx context.Context, db *sql.DB) ([]model.Follow, error) {
	q := "SELECT `id`, `following_user_id`, `followed_user_id`, `created_at`, `updated_at` FROM `follows`"
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []model.Follow{}
	for rows.Next() {
		var f model.Follow
		if err := rows.Scan(&f.ID, &f.FollowingUserID, &f.FollowedUserID, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, f)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func FindAllUserReports(ctx context.Context, db *sql.DB) ([]model.UserReport, error) {
	q := "SELECT `id`, `user_name`, `detail`, `created_at`, `updated_at` FROM `user_reports`"
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []model.UserReport{}
	for rows.Next() {
		var u model.UserReport
		if err := rows.Scan(&u.ID, &u.UserName, &u.Detail, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, u)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func FindAllPostingReports(ctx context.Context, db *sql.DB) ([]model.PostingReport, error) {
	q := "SELECT `id`, `posting_id`, `detail`, `created_at`, `updated_at` FROM `posting_reports`"
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []model.PostingReport{}
	for rows.Next() {
		var u model.PostingReport
		if err := rows.Scan(&u.ID, &u.PostingID, &u.Detail, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, u)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func CreatePasswordReset(db *sql.DB, expiresAt time.Time) error {
	q := "INSERT INTO `password_resets` (`user_id`, `password_reset_email_count`, `password_reset_key`, `password_reset_key_expires_at`) VALUES (?, ?, ?, ?)"
	_, e := db.Exec(q, dummy.User1.ID, dummy.PasswordReset1.PasswordResetEmailCount, dummy.PasswordReset1.PasswordResetKey, expiresAt)
	if e != nil {
		return e
	}
	return nil
}

func UpdatePasswordResetExpiresAt(db *sql.DB, expiresAt time.Time) error {
	q := "UPDATE `password_resets` SET `password_reset_key_expires_at` = ?"
	_, e := db.Exec(q, expiresAt)
	if e != nil {
		return e
	}
	return nil
}

func UpdatePasswordResetEmailCount(db *sql.DB) error {
	q := "UPDATE `password_resets` SET `password_reset_email_count` = 10"
	_, e := db.Exec(q)
	if e != nil {
		return e
	}
	return nil
}

func DeleteAllTableData(db *sql.DB, table string) error {
	q := fmt.Sprintf("DELETE FROM `%s`", table)
	_, e := db.Exec(q)
	if e != nil {
		return e
	}
	return nil
}

func UpdateNow(db *sql.DB, table string) error {
	q := fmt.Sprintf("UPDATE `%s` SET `created_at` = ?, `updated_at` = ?", table)
	_, e := db.Exec(q, lib.NowFunc(), lib.NowFunc())
	if e != nil {
		return e
	}
	return nil
}
