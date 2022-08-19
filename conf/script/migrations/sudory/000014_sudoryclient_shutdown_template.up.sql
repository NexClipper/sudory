INSERT INTO `template` (`uuid`, `name`, `summary`, `origin`, `created`) VALUES ('99990000000000000000000000000001', 'sudoryclient_shutdown', 'class=\'sudoryclient\' resource=\'shutdown\' verb=\'\' group_version=\'\'', 'predefined', NOW());
INSERT INTO `template_command` (`uuid`, `name`, `summary`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`, `created`) VALUES ('99990000000000000000000000000001', 'sudoryclient_shutdown_0', 'class=\'sudoryclient\' resource=\'shutdown\' verb=\'\' group_version=\'\'', '99990000000000000000000000000001', '0', 'sudoryclient.shutdown', '{}', NULL, NOW());

INSERT INTO `template_recipe` ( `id`, `method`, `args`, `name`, `summary`, `created`) VALUES ('133', 'sudoryclient.shutdown', '{}', 'sudoryclient-shutdown', '', NOW());
