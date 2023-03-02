CREATE TABLE `qr_menu`.`collections` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `collection_id` BIGINT UNSIGNED NULL DEFAULT NULL,
    `store_id` BIGINT UNSIGNED NOT NULL,
    `user_app_token_id` BIGINT UNSIGNED NULL DEFAULT NULL,
    `description` TEXT NULL DEFAULT NULL,
    `alias` TEXT NULL DEFAULT NULL,
    `name` TEXT NOT NULL,
    `image` TEXT NULL DEFAULT NULL,
    `gateway` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP on update CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX (`collection_id`),
    INDEX (`user_app_token_id`),
    INDEX (`gateway`),
    INDEX (`store_id`)
) ENGINE = InnoDB CHARSET = utf8mb4 COLLATE utf8mb4_general_ci;

ALTER TABLE
    `qr_menu`.`collections`
ADD
    CONSTRAINT `store_collection_unique` UNIQUE KEY(store_id, collection_id);