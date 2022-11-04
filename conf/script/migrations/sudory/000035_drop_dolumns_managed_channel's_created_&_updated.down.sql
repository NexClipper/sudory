-- managed_channel_format
ALTER TABLE `managed_channel_format`
	ADD COLUMN IF NOT EXISTS `created` DATETIME NOT NULL DEFAULT NOW();
ALTER TABLE `managed_channel_format`
	ADD COLUMN IF NOT EXISTS `updated` DATETIME NULL;

-- managed_channel_notifier_console
ALTER TABLE `managed_channel_notifier_console`
	ADD COLUMN IF NOT EXISTS `created` DATETIME NOT NULL DEFAULT NOW();
ALTER TABLE `managed_channel_notifier_console`
	ADD COLUMN IF NOT EXISTS `updated` DATETIME NULL;

-- managed_channel_notifier_edge
ALTER TABLE `managed_channel_notifier_edge`
	ADD COLUMN IF NOT EXISTS `created` DATETIME NOT NULL DEFAULT NOW();
ALTER TABLE `managed_channel_notifier_edge`
	ADD COLUMN IF NOT EXISTS `updated` DATETIME NULL;

-- managed_channel_notifier_rabbitmq
ALTER TABLE `managed_channel_notifier_rabbitmq`
	ADD COLUMN IF NOT EXISTS `created` DATETIME NOT NULL DEFAULT NOW();
ALTER TABLE `managed_channel_notifier_rabbitmq`
	ADD COLUMN IF NOT EXISTS `updated` DATETIME NULL;

-- managed_channel_notifier_slackhook
ALTER TABLE `managed_channel_notifier_slackhook`
	ADD COLUMN IF NOT EXISTS `created` DATETIME NOT NULL DEFAULT NOW();
ALTER TABLE `managed_channel_notifier_slackhook`
	ADD COLUMN IF NOT EXISTS `updated` DATETIME NULL;

-- managed_channel_notifier_webhook
ALTER TABLE `managed_channel_notifier_webhook`
	ADD COLUMN IF NOT EXISTS `created` DATETIME NOT NULL DEFAULT NOW();
ALTER TABLE `managed_channel_notifier_webhook`
	ADD COLUMN IF NOT EXISTS `updated` DATETIME NULL;

-- managed_channel_status_option
ALTER TABLE `managed_channel_status_option`
	ADD COLUMN IF NOT EXISTS `created` DATETIME NOT NULL DEFAULT NOW();
ALTER TABLE `managed_channel_status_option`
	ADD COLUMN IF NOT EXISTS `updated` DATETIME NULL;

