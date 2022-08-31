-- service
ALTER TABLE `service`
	CHANGE COLUMN `created` `created` DATETIME NOT NULL AFTER `message`;
ALTER TABLE `service_step`
	CHANGE COLUMN `created` `created` DATETIME NOT NULL AFTER `ended`;

-- cluster_information
ALTER TABLE `cluster_information`
	CHANGE COLUMN `polling_offset` `polling_offset` DATETIME DEFAULT NULL AFTER `polling_count`;
