-- service
ALTER TABLE `service`
	CHANGE COLUMN `created` `created` DATETIME(6) NOT NULL AFTER `message`;
ALTER TABLE `service_step`
	CHANGE COLUMN `created` `created` DATETIME(6) NOT NULL AFTER `ended`;

-- cluster_information
ALTER TABLE `cluster_information`
	CHANGE COLUMN `polling_offset` `polling_offset` DATETIME(6) DEFAULT NULL AFTER `polling_count`;
