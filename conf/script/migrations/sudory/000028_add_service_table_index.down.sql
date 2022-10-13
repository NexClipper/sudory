-- service

ALTER TABLE `service`
	DROP INDEX IF EXISTS `uuid`;

ALTER TABLE `service`
	DROP INDEX IF EXISTS `created`;
	
ALTER TABLE `service`
	DROP INDEX IF EXISTS `cluster_uuid`;
	
-- service_step

ALTER TABLE `service_step`
	DROP INDEX IF EXISTS `uuid`;

ALTER TABLE `service_step`
	DROP INDEX IF EXISTS `created`;
	
ALTER TABLE `service_step`
	DROP INDEX IF EXISTS `cluster_uuid`;	

-- service_result

ALTER TABLE `service_result`
	DROP INDEX IF EXISTS `uuid`;
	
ALTER TABLE `service_result`
	DROP INDEX IF EXISTS `cluster_uuid`;	