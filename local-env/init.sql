DROP DATABASE IF EXISTS app;
CREATE DATABASE app;
USE app;
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