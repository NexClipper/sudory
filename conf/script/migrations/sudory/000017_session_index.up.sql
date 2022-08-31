

ALTER TABLE `session`
    DROP INDEX IF EXISTS `cluster_expiration_time`;

ALTER TABLE `session`
    DROP INDEX IF EXISTS `cluster_uuid`;

ALTER TABLE `session`
	ADD INDEX IF NOT EXISTS `cluster_expiration_time` (`cluster_uuid`, `expiration_time`);

ALTER TABLE `session`
	ADD INDEX IF NOT EXISTS `cluster_uuid` (`cluster_uuid`, `uuid`);