DROP TABLE IF EXISTS `tenant_clusters`;
DROP TABLE IF EXISTS `tenant_channels`;
DROP TABLE IF EXISTS `tenant`;

-- alter channel tables
ALTER TABLE `managed_channel_filter`
	CHANGE COLUMN IF EXISTS `created` `created` DATETIME NOT NULL;

ALTER TABLE `managed_channel_format`
	CHANGE COLUMN IF EXISTS `created` `created` DATETIME NOT NULL;

ALTER TABLE `managed_channel_notifier_console`
	CHANGE COLUMN IF EXISTS `created` `created` DATETIME NOT NULL;

ALTER TABLE `managed_channel_notifier_edge`
	CHANGE COLUMN IF EXISTS `created` `created` DATETIME NOT NULL;

ALTER TABLE `managed_channel_notifier_rabbitmq`
	CHANGE COLUMN IF EXISTS `created` `created` DATETIME NOT NULL;

ALTER TABLE `managed_channel_notifier_slackhook`
	CHANGE COLUMN IF EXISTS `created` `created` DATETIME NOT NULL;

ALTER TABLE `managed_channel_notifier_webhook`
	CHANGE COLUMN IF EXISTS `created` `created` DATETIME NOT NULL;

ALTER TABLE `managed_channel_status_option`
	CHANGE COLUMN IF EXISTS `created` `created` DATETIME NOT NULL;

