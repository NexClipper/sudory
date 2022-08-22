INSERT INTO `template` (`uuid`, `name`, `summary`, `origin`, `created`) VALUES ('99990000000000000000000000000001', 'sudory_client_pod_delete', 'class=\'sudory\' resource=\'client_pod\' verb=\'delete\' group_version=\'\'', 'predefined', NOW());
INSERT INTO `template_command` (`uuid`, `name`, `summary`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`, `created`) VALUES ('99990000000000000000000000000001', 'sudory_client_pod_delete_0', 'class=\'sudory\' resource=\'client_pod\' verb=\'delete\' group_version=\'\'', '99990000000000000000000000000001', '0', 'sudory.client_pod.delete', '{}', NULL, NOW());

INSERT INTO `template_recipe` ( `id`, `method`, `args`, `name`, `summary`, `created`) VALUES ('133', 'sudory.client_pod.delete', '{}', 'sudory-client_pod-delete', '', NOW());
