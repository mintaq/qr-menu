CREATE TABLE `qr_menu`.`apps` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `app_name` VARCHAR(255) NOT NULL,
    `api_key` MEDIUMTEXT NOT NULL,
    `secret_key` MEDIUMTEXT NOT NULL,
    `scopes` MEDIUMTEXT NULL,
    `redirect_url` VARCHAR(255) NULL,
    `gateway` ENUM('sapo', 'kiotviet') NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP on update CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE (`app_name`)
) ENGINE = InnoDB CHARSET = utf8mb4 COLLATE utf8mb4_general_ci;