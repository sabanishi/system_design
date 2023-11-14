-- Table for tasks
DROP TABLE IF EXISTS `tasks`;
DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `ownership`;

CREATE TABLE `users`(
    `id`            bigint(20) NOT NULL AUTO_INCREMENT,
    `name`          varchar(50) NOT NULL UNIQUE,
    `password`      binary(32) NOT NULL,
    `updated_at`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_at`    datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `is_deleted`    boolean NOT NULL DEFAULT b'0',
    PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8mb4;

CREATE TABLE `tasks`
(
    `id`         bigint(20) NOT NULL AUTO_INCREMENT,
    `title`      varchar(50) NOT NULL,
    `is_done`    boolean     NOT NULL DEFAULT b'0',
    `created_at` datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `description`   varchar(256) NOT NULL,
    PRIMARY KEY (`id`)
) DEFAULT CHARSET=utf8mb4;

CREATE TABLE `ownership`(
    `user_id` bigint(20) NOT NULL,
    `task_id` bigint(20) NOT NULL,
    PRIMARY KEY (`user_id`, `task_id`)
) DEFAULT CHARSET=utf8mb4;
