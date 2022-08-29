
-- managed_channel_filter
ALTER TABLE `managed_channel_filter`
	DROP COLUMN IF EXISTS `created`;
ALTER TABLE `managed_channel_filter`
	DROP COLUMN IF EXISTS `updated`;

-- managed_channel_format
ALTER TABLE `managed_channel_format`
	DROP COLUMN IF EXISTS `created`;
ALTER TABLE `managed_channel_format`
	DROP COLUMN IF EXISTS `updated`;

-- managed_channel_notifier_console
ALTER TABLE `managed_channel_notifier_console`
	DROP COLUMN IF EXISTS `created`;
ALTER TABLE `managed_channel_notifier_console`
	DROP COLUMN IF EXISTS `updated`;

-- managed_channel_notifier_edge
ALTER TABLE `managed_channel_notifier_edge`
	DROP COLUMN IF EXISTS `created`;
ALTER TABLE `managed_channel_notifier_edge`
	DROP COLUMN IF EXISTS `updated`;

-- managed_channel_notifier_rabbitmq
ALTER TABLE `managed_channel_notifier_rabbitmq`
	DROP COLUMN IF EXISTS `created`;
ALTER TABLE `managed_channel_notifier_rabbitmq`
	DROP COLUMN IF EXISTS `updated`;

-- managed_channel_notifier_slackhook
ALTER TABLE `managed_channel_notifier_slackhook`
	DROP COLUMN IF EXISTS `created`;
ALTER TABLE `managed_channel_notifier_slackhook`
	DROP COLUMN IF EXISTS `updated`;

-- managed_channel_notifier_webhook
ALTER TABLE `managed_channel_notifier_webhook`
	DROP COLUMN IF EXISTS `created`;
ALTER TABLE `managed_channel_notifier_webhook`
	DROP COLUMN IF EXISTS `updated`;

-- managed_channel_status_option
ALTER TABLE `managed_channel_status_option`
	DROP COLUMN IF EXISTS `created`;
ALTER TABLE `managed_channel_status_option`
	DROP COLUMN IF EXISTS `updated`;
