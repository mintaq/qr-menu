CREATE TABLE `qr_menu`.`user_app_tokens` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `app_id` BIGINT UNSIGNED NOT NULL,
    `store_domain` VARCHAR(255) NOT NULL,
    `access_token` MEDIUMTEXT NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP on update CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX (`user_id`),
    INDEX (`app_id`),
    INDEX (`store_domain`)
) ENGINE = InnoDB CHARSET = utf8mb4 COLLATE utf8mb4_general_ci;

ALTER TABLE `qr_menu`.`user_app_tokens` ADD CONSTRAINT `user_app_unique` UNIQUE KEY(user_id, app_id)