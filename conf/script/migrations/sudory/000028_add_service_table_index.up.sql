-- service

ALTER TABLE `service`
	ADD INDEX IF NOT EXISTS `uuid` (`uuid`);

ALTER TABLE `service`
	ADD INDEX IF NOT EXISTS `created` (`created`);
	
ALTER TABLE `service`
	ADD INDEX IF NOT EXISTS `cluster_uuid` (`cluster_uuid`);

-- service_step

ALTER TABLE `service_step`
	ADD INDEX IF NOT EXISTS `uuid` (`uuid`);

ALTER TABLE `service_step`
	ADD INDEX IF NOT EXISTS `created` (`created`);
	
ALTER TABLE `service_step`
	ADD INDEX IF NOT EXISTS `cluster_uuid` (`cluster_uuid`);

-- service_result

ALTER TABLE `service_result`
	ADD INDEX IF NOT EXISTS `uuid` (`uuid`);
	
ALTER TABLE `service_result`
	ADD INDEX IF NOT EXISTS `cluster_uuid` (`cluster_uuid`);