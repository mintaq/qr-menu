CREATE TABLE `qr_menu`.`menus` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `store_id` BIGINT UNSIGNED NOT NULL,
    `name` VARCHAR(255) NULL,
    `role` VARCHAR(25) NULL DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP on update CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX (`store_id`),
    INDEX (`name`),
    INDEX (`role`)
) ENGINE = InnoDB CHARSET = utf8mb4 COLLATE utf8mb4_general_ci;

ALTER TABLE
    `qr_menu`.`menus`
ADD
    CONSTRAINT `store_name_unique` UNIQUE KEY(store_id, name);