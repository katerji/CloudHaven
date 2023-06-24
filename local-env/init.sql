DROP DATABASE IF EXISTS APP;
CREATE DATABASE APP;
USE APP;
CREATE TABLE `user`
(
    `id`         int          NOT NULL AUTO_INCREMENT,
    `email`      varchar(255) NOT NULL,
    `name`       varchar(255) NOT NULL,
    `password`   varchar(255) NOT NULL,
    `created_on` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_on` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `file`
(
    `id`           int          NOT NULL AUTO_INCREMENT,
    `name`         varchar(255) NOT NULL,
    `size`         bigint       NOT NULL,
    `content_type` varchar(255) NOT NULL,
    `owner_id`     int          NOT NULL,
    `created_on`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_on`   datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`,`owner_id`),
    KEY            `owner_id` (`owner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `file_share`
(
    `id`         int      NOT NULL AUTO_INCREMENT,
    `file_id`    int      NOT NULL,
    `url`        text     NOT NULL,
    `open_rate`  int      NOT NULL DEFAULT '0',
    `expires_at` int      NOT NULL,
    `created_on` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `file_id` (`file_id`,`expires_at`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4;