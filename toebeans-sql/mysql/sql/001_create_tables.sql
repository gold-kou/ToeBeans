CREATE TABLE `users` (
    `name` VARCHAR(255) PRIMARY KEY,
    `email` VARCHAR(255) UNIQUE NOT NULL,
    `password` CHAR(128) NOT NULL,
    `icon` VARCHAR(255) NOT NULL DEFAULT 'UNKNOWN',
    `self_introduction` VARCHAR(255) NOT NULL DEFAULT 'UNKNOWN',
    `posting_count` INT NOT NULL DEFAULT 0,
    `like_count` INT NOT NULL DEFAULT 0,
    `liked_count` INT NOT NULL DEFAULT 0,
    `follow_count` INT NOT NULL DEFAULT 0,
    `followed_count` INT NOT NULL DEFAULT 0,
    `activation_key` VARCHAR(255) NOT NULL,
    `email_verified` BOOLEAN NOT NULL DEFAULT FALSE,
    `password_reset_email_count` INT NOT NULL DEFAULT 0,
    `password_reset_key` VARCHAR(255) NOT NULL DEFAULT 'UNKNOWN',
    `password_reset_key_expires_at` DATETIME NOT NULL DEFAULT '1000-01-01 00:00:00',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
    INDEX idx_users_name(name)
);

CREATE TABLE `postings` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `user_name` VARCHAR(255) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `image_url` VARCHAR(255) NOT NULL,
    `liked_count` INT NOT NULL DEFAULT 0,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `postings_user_name` FOREIGN KEY (`user_name`) REFERENCES `users` (`name`),
    INDEX idx_postings_user_name(user_name)
);

CREATE TABLE `likes` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `user_name` VARCHAR(255) NOT NULL,
    `posting_id` INT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `likes_user_name` FOREIGN KEY (`user_name`) REFERENCES `users` (`name`),
    CONSTRAINT `likes_posting_id` FOREIGN KEY (`posting_id`) REFERENCES `postings` (`id`),
    UNIQUE `uk_user_name_posting_id` (`user_name`, `posting_id`)
);

CREATE TABLE `comments` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `user_name` VARCHAR(255) NOT NULL,
    `posting_id` INT NOT NULL,
    `comment` VARCHAR(255) NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `comments_user_name` FOREIGN KEY (`user_name`) REFERENCES `users` (`name`),
    CONSTRAINT `comments_posting_id` FOREIGN KEY (`posting_id`) REFERENCES `postings` (`id`),
    INDEX idx_comments_posting_id(posting_id)
);

CREATE TABLE `follows` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `following_user_name` VARCHAR(255) NOT NULL,
    `followed_user_name` VARCHAR(255) NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `follows_following_user_name` FOREIGN KEY (`following_user_name`) REFERENCES `users` (`name`),
    CONSTRAINT `follows_followed_user_name` FOREIGN KEY (`followed_user_name`) REFERENCES `users` (`name`)
);

CREATE TABLE `notifications` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `visitor_name` VARCHAR(255) NOT NULL,
    `visited_name` VARCHAR(255) NOT NULL,
    `posting_id` INT NOT NULL,
    `comment_id` INT NOT NULL,
    `action` ENUM('like', 'comment', 'follow') NOT NULL,
    `checked` BOOLEAN NOT NULL DEFAULT FALSE,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `notifications_visitor_name` FOREIGN KEY (`visitor_name`) REFERENCES `users` (`name`),
    CONSTRAINT `notifications_visited_name` FOREIGN KEY (`visited_name`) REFERENCES `users` (`name`),
    CONSTRAINT `notifications_posting_id` FOREIGN KEY (`posting_id`) REFERENCES `postings` (`id`),
    CONSTRAINT `notifications_comment_id` FOREIGN KEY (`comment_id`) REFERENCES `comments` (`id`),
    INDEX idx_notifications_visited_name(visited_name)
);
