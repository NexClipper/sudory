-- service.priority
ALTER TABLE `service`
	ADD COLUMN IF NOT EXISTS `priority` TINYINT UNSIGNED NOT NULL DEFAULT '0' AFTER `step_count`;