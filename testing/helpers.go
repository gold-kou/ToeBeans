package testing

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/gold-kou/ToeBeans/app/lib"

	"github.com/gold-kou/ToeBeans/app/domain/model"

	"github.com/gold-kou/ToeBeans/app/adapter/mysql"
)

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
	db.Close()
}

func FindAllUsers(ctx context.Context, db *sql.DB) ([]model.User, error) {
	q := "SELECT `name`, `email`, `password`, `icon`, `self_introduction`, `posting_count`, `like_count`, `liked_count`, `follow_count`, `followed_count`, `activation_key`, `email_verified`, `password_reset_email_count`, `password_reset_key`, `password_reset_key_expires_at`, `created_at`, `updated_at` FROM `users`"
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []model.User{}
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.Name, &u.Email, &u.Password, &u.Icon, &u.SelfIntroduction, &u.PostingCount, &u.LikeCount, &u.LikedCount, &u.FollowCount, &u.FollowedCount, &u.ActivationKey, &u.EmailVerified, &u.PasswordResetEmailCount, &u.PasswordResetKey, &u.PasswordResetKeyExpiresAt, &u.CreatedAt, &u.UpdatedAt); err != nil {
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
