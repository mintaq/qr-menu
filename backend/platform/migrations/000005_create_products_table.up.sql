CREATE TABLE `qr_menu`.`products` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `store_id` BIGINT UNSIGNED NOT NULL,
    `user_app_token_id` BIGINT UNSIGNED NULL DEFAULT NULL,
    `content` MEDIUMTEXT NULL DEFAULT NULL COMMENT 'The HTML of product',
    `summary` MEDIUMTEXT NULL DEFAULT NULL COMMENT 'The short description of product',
    `created_on` TIMESTAMP NULL,
    `alias` TEXT NULL COMMENT 'The unique string represents the product',
    `product_id` BIGINT UNSIGNED NULL DEFAULT NULL,
    `images` JSON NULL DEFAULT '[]',
    `options` JSON NULL DEFAULT '[]',
    `product_type` TEXT NULL DEFAULT NULL,
    `product_status` VARCHAR(50) NULL DEFAULT 'active',
    `price` FLOAT NULL DEFAULT NULL,
    `published_on` TIMESTAMP NULL DEFAULT NULL,
    `tags` MEDIUMTEXT NULL DEFAULT NULL,
    `product_name` TEXT NOT NULL,
    `modified_on` TIMESTAMP NULL DEFAULT NULL,
    `variants` JSON NULL DEFAULT '[]',
    `vendor` TEXT NULL DEFAULT NULL,
    `is_charge_tax` BOOLEAN NOT NULL DEFAULT FALSE,
    `gateway` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP on update CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX (`store_id`),
    INDEX (`product_id`),
    INDEX (`user_app_token_id`),
    INDEX (`gateway`),
    INDEX (`product_type`)
) ENGINE = InnoDB CHARSET = utf8mb4 COLLATE utf8mb4_general_ci;

ALTER TABLE
    `qr_menu`.`products`
ADD
    CONSTRAINT `store_product_unique` UNIQUE KEY(store_id, product_id);