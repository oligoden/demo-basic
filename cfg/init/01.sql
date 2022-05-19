CREATE DATABASE IF NOT EXISTS `test`;
GRANT ALL ON `test`.* TO 'user'@'%';

CREATE TABLE IF NOT EXISTS `routings`
(
    `domain` varchar(255) DEFAULT NULL,
    `path` varchar(255) DEFAULT NULL,
    `url` varchar(255) DEFAULT NULL,
    `reset_cors` boolean,
    `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
    `uc` varchar(255) DEFAULT NULL,
    `owner_id` int(10) unsigned DEFAULT NULL,
    `perms` varchar(255) DEFAULT NULL,
    `hash` varchar(255) DEFAULT NULL,
    PRIMARY KEY (`id`), UNIQUE KEY `uc` (`uc`)
);
INSERT INTO `routings`
(`uc`, `domain`, `path`, `url`, `reset_cors`, `owner_id`, `perms`, `hash`)
VALUES
('a', 'api.test.com:8080', '/', 'swagger:8080', false, 1, ':::', 'abc');