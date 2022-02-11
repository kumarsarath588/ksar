CREATE TABLE IF NOT EXISTS `tabsquare`.`customers` (
    `uuid` VARCHAR(255) NOT NULL,
    `customer_name` VARCHAR(2048) NOT NULL,
    `country` VARCHAR(4096) NOT NULL,
    PRIMARY KEY (`uuid`)
);