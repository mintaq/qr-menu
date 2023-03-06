CREATE TABLE `qr_menu`.`collects` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `store_id` BIGINT UNSIGNED NOT NULL,
    `user_app_token_id` BIGINT UNSIGNED NULL DEFAULT NULL,
    `collection_id` BIGINT UNSIGNED NOT NULL,
    `product_id` BIGINT UNSIGNED NOT NULL,
    `position` INT NULL DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP on update CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX (`store_id`),
    INDEX (`collection_id`),
    INDEX (`user_app_token_id`),
    INDEX (`product_id`)
) ENGINE = InnoDB CHARSET = utf8mb4 COLLATE utf8mb4_general_ci;

ALTER TABLE
    `qr_menu`.`collects`
ADD
    CONSTRAINT `store_collection_product_unique` UNIQUE KEY(store_id, collection_id, product_id);