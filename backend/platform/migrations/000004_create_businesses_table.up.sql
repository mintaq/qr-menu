CREATE TABLE `qr_menu`.`businesses` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `store_name` VARCHAR(255) NOT NULL,
    `store_url` VARCHAR(255) NOT NULL,
    `country` VARCHAR(255) NOT NULL,
    `city` VARCHAR(255) NOT NULL,
    `address` MEDIUMTEXT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP on update CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX (`user_id`)
) ENGINE = InnoDB CHARSET = utf8mb4 COLLATE utf8mb4_general_ci;
