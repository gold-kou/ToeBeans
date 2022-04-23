CREATE TABLE `users` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) UNIQUE NOT NULL,
    `email` VARCHAR(255) UNIQUE NOT NULL,
    `password` CHAR(128) NOT NULL,
    `icon` VARCHAR(255) NOT NULL DEFAULT 'UNKNOWN',
    `self_introduction` VARCHAR(255) NOT NULL DEFAULT 'UNKNOWN',
    `activation_key` VARCHAR(255) NOT NULL,
    `email_verified` BOOLEAN NOT NULL DEFAULT FALSE,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
    INDEX idx_users_name(name)
);

CREATE TABLE `password_resets` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `user_id` INT UNIQUE NOT NULL,
    `password_reset_email_count` INT NOT NULL DEFAULT 0,
    `password_reset_key` VARCHAR(255) NOT NULL DEFAULT 'UNKNOWN',
    `password_reset_key_expires_at` DATETIME NOT NULL DEFAULT '1000-01-01 00:00:00',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `password_resets_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
);

CREATE TABLE `postings` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `user_id` INT NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `image_url` VARCHAR(255) NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `postings_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
    INDEX idx_postings_user_id(user_id)
);

CREATE TABLE `likes` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `user_id` INT NOT NULL,
    `posting_id` INT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `likes_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
    CONSTRAINT `likes_posting_id` FOREIGN KEY (`posting_id`) REFERENCES `postings` (`id`),
    UNIQUE `uk_user_id_posting_id` (`user_id`, `posting_id`)
);

CREATE TABLE `comments` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `user_id` INT NOT NULL,
    `posting_id` INT NOT NULL,
    `comment` VARCHAR(255) NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `comments_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
    CONSTRAINT `comments_posting_id` FOREIGN KEY (`posting_id`) REFERENCES `postings` (`id`),
    INDEX idx_comments_posting_id(posting_id)
);

CREATE TABLE `follows` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `following_user_id` INT NOT NULL,
    `followed_user_id` INT NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `follows_following_user_id` FOREIGN KEY (`following_user_id`) REFERENCES `users` (`id`),
    CONSTRAINT `follows_followed_user_id` FOREIGN KEY (`followed_user_id`) REFERENCES `users` (`id`),
    UNIQUE `uk_following_user_id_followed_user_id` (`following_user_id`, `followed_user_id`)
);

CREATE TABLE `notifications` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `visitor_user_id` INT NOT NULL,
    `visited_user_id` INT NOT NULL,
    `action` ENUM('like', 'comment', 'follow') NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `notifications_visitor_user_id` FOREIGN KEY (`visitor_user_id`) REFERENCES `users` (`id`),
    CONSTRAINT `notifications_visited_user_id` FOREIGN KEY (`visited_user_id`) REFERENCES `users` (`id`),
    INDEX idx_notifications_visited_user_id(visited_user_id)
);

CREATE TABLE `user_reports` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `user_name` VARCHAR(255) NOT NULL,
    `detail` VARCHAR(3000) NOT NULL DEFAULT 'UNKNOWN',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `user_reports_user_name` FOREIGN KEY (`user_name`) REFERENCES `users` (`name`),
    INDEX idx_user_reports_created_at(created_at)
);

CREATE TABLE `posting_reports` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `posting_id` INT NOT NULL,
    `detail` VARCHAR(3000) NOT NULL DEFAULT 'UNKNOWN',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT `posting_reports_posting_id` FOREIGN KEY (`posting_id`) REFERENCES `postings` (`id`),
    INDEX idx_posting_reports_created_at(created_at)
);