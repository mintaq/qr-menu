CREATE TABLE `qr_menu`.`orders` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `store_id` BIGINT UNSIGNED NOT NULL,
    `billing_address` JSON NULL DEFAULT NULL,
    `browser_ip` VARCHAR(255) NULL DEFAULT NULL,
    `buyer_accepts_marketing` TINYINT(1) NULL DEFAULT NULL,
    `cancel_reason` TEXT NULL DEFAULT NULL,
    `cancelled_on` TIMESTAMP NULL DEFAULT NULL,
    `cart_token` TEXT NULL DEFAULT NULL,
    `client_details` JSON NULL DEFAULT NULL,
    `closed_on` TIMESTAMP NULL DEFAULT NULL,
    `currency` VARCHAR(5) NOT NULL,
    `customer` JSON NULL DEFAULT NULL,
    `discount_codes` JSON NULL DEFAULT NULL,
    `email` VARCHAR(255) NULL DEFAULT NULL,
    `financial_status` VARCHAR(255) NULL DEFAULT NULL,
    `status` VARCHAR(25) NULL DEFAULT NULL,
    `fulfillments` JSON NULL DEFAULT NULL,
    `fulfillment_status` VARCHAR(25) NULL DEFAULT NULL,
    `tags` TEXT NULL DEFAULT NULL,
    `landing_site` VARCHAR(255) NULL DEFAULT NULL,
    `line_items` JSON NULL DEFAULT NULL,
    `name` VARCHAR(255) NULL DEFAULT NULL,
    `note` TEXT NULL DEFAULT NULL,
    `note_attributes` JSON NULL DEFAULT NULL,
    `number` BIGINT NULL DEFAULT NULL,
    `order_number` BIGINT NULL DEFAULT NULL,
    `payment_gate_way_names` JSON NULL DEFAULT NULL,
    `processed_on` TIMESTAMP NULL DEFAULT NULL,
    `processing_method` VARCHAR(50) NULL DEFAULT NULL,
    `referring_site` VARCHAR(255) NULL DEFAULT NULL,
    `refunds` TEXT NULL DEFAULT NULL,
    `shipping_address` JSON NULL DEFAULT NULL,
    `shipping_lines` JSON NULL DEFAULT NULL,
    `source_name` TEXT NULL DEFAULT NULL,
    `token` TEXT NULL DEFAULT NULL,
    `total_discount` FLOAT NULL DEFAULT NULL,
    `total_line_items_price` FLOAT NULL DEFAULT NULL,
    `total_price` FLOAT NULL DEFAULT NULL,
    `total_weight` INT NULL DEFAULT NULL,
    `modified_on` TIMESTAMP NULL DEFAULT NULL,
    `gateway` VARCHAR(255) NULL DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP on update CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX (`store_id`),
    INDEX (`cart_token`)
) ENGINE = InnoDB CHARSET = utf8mb4 COLLATE utf8mb4_general_ci;