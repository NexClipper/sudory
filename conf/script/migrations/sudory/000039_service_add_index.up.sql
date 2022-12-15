ALTER TABLE `service`
	ADD INDEX IF NOT EXISTS `template_uuid` (`cluster_uuid`, `template_uuid`);
