
-- drop newist table
DROP TABLE IF EXISTS `service`;
DROP TABLE IF EXISTS `service_step`;
DROP TABLE IF EXISTS `service_result`;

-- restore old tables
RENAME TABLE `old_service` TO `service`;
RENAME TABLE `old_service_status` TO `service_status`;
RENAME TABLE `old_service_step` TO `service_step`;
RENAME TABLE `old_service_step_status` TO `service_step_status`;
RENAME TABLE `old_service_result` TO `service_result`;
