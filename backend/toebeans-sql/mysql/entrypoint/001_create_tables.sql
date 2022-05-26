CREATE TABLE `users` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'サロゲートキー',
    `name` VARCHAR(255) UNIQUE NOT NULL COMMENT 'ユーザ名',
    `email` VARCHAR(255) UNIQUE NOT NULL COMMENT 'メールアドレス',
    `password` CHAR(128) NOT NULL COMMENT 'ハッシュ済みパスワード',
    `icon` VARCHAR(255) NOT NULL DEFAULT 'UNKNOWN' COMMENT 'アイコン',
    `self_introduction` VARCHAR(255) NOT NULL DEFAULT 'UNKNOWN' COMMENT '自己紹介文',
    `activation_key` VARCHAR(255) NOT NULL COMMENT 'アクティベーションキー。UUIDで、ユーザ登録メール内のリンクに付与される。',
    `email_verified` BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'メール本人確認が済んでいるかどうか',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    INDEX idx_users_name(name)
)COMMENT 'ユーザテーブル';

CREATE TABLE `password_resets` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'サロゲートキー',
    `user_id` INT UNIQUE NOT NULL COMMENT 'ユーザID',
    `password_reset_email_count` INT NOT NULL DEFAULT 0 COMMENT 'パスワードリセットした回数。一日あたりのリセット回数を制限する用途で用意したがバッチは未実装。',
    `password_reset_key` VARCHAR(255) NOT NULL DEFAULT 'UNKNOWN' COMMENT 'パスワードリセットキー',
    `password_reset_key_expires_at` DATETIME NOT NULL DEFAULT '1000-01-01 00:00:00' COMMENT 'パスワードリセットキーの有効期限',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    CONSTRAINT `password_resets_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
)COMMENT 'パスワードリセットテーブル';

CREATE TABLE `postings` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'サロゲートキー',
    `user_id` INT NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `image_url` VARCHAR(255) NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    CONSTRAINT `postings_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
    INDEX idx_postings_user_id(user_id)
)COMMENT '投稿テーブル';

CREATE TABLE `likes` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'サロゲートキー',
    `user_id` INT NOT NULL,
    `posting_id` INT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    CONSTRAINT `likes_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
    CONSTRAINT `likes_posting_id` FOREIGN KEY (`posting_id`) REFERENCES `postings` (`id`),
    UNIQUE `uk_user_id_posting_id` (`user_id`, `posting_id`)
)COMMENT 'いいねテーブル';

CREATE TABLE `comments` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'サロゲートキー',
    `user_id` INT NOT NULL,
    `posting_id` INT NOT NULL,
    `comment` VARCHAR(255) NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    CONSTRAINT `comments_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
    CONSTRAINT `comments_posting_id` FOREIGN KEY (`posting_id`) REFERENCES `postings` (`id`),
    INDEX idx_comments_posting_id(posting_id)
)COMMENT 'コメントテーブル';

CREATE TABLE `follows` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'サロゲートキー',
    `following_user_id` INT NOT NULL,
    `followed_user_id` INT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    CONSTRAINT `follows_following_user_id` FOREIGN KEY (`following_user_id`) REFERENCES `users` (`id`),
    CONSTRAINT `follows_followed_user_id` FOREIGN KEY (`followed_user_id`) REFERENCES `users` (`id`),
    UNIQUE `uk_following_user_id_followed_user_id` (`following_user_id`, `followed_user_id`)
)COMMENT 'フォローテーブル';

CREATE TABLE `notifications` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'サロゲートキー',
    `visitor_user_id` INT NOT NULL,
    `visited_user_id` INT NOT NULL,
    `action` ENUM('like', 'comment', 'follow') NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    CONSTRAINT `notifications_visitor_user_id` FOREIGN KEY (`visitor_user_id`) REFERENCES `users` (`id`),
    CONSTRAINT `notifications_visited_user_id` FOREIGN KEY (`visited_user_id`) REFERENCES `users` (`id`),
    INDEX idx_notifications_visited_user_id(visited_user_id)
)COMMENT '通知テーブル';

CREATE TABLE `user_reports` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'サロゲートキー',
    `user_name` VARCHAR(255) NOT NULL,
    `detail` VARCHAR(3000) NOT NULL DEFAULT 'UNKNOWN',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    CONSTRAINT `user_reports_user_name` FOREIGN KEY (`user_name`) REFERENCES `users` (`name`),
    INDEX idx_user_reports_created_at(created_at)
)COMMENT 'ユーザレポートテーブル';

CREATE TABLE `posting_reports` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT 'サロゲートキー',
    `posting_id` INT NOT NULL,
    `detail` VARCHAR(3000) NOT NULL DEFAULT 'UNKNOWN',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',
    CONSTRAINT `posting_reports_posting_id` FOREIGN KEY (`posting_id`) REFERENCES `postings` (`id`),
    INDEX idx_posting_reports_created_at(created_at)
)COMMENT '投稿レポートテーブル';