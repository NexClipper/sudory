use sudory;

-- template
-- k8s.group_version: v1
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000001', 'kubernetes_pods_get', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000002', 'kubernetes_pods_list', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000003', 'kubernetes_namespaces_get', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000004', 'kubernetes_namespaces_list', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000005', 'kubernetes_configmaps_get', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000006', 'kubernetes_configmaps_list', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000007', 'kubernetes_events_get', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000008', 'kubernetes_events_list', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000009', 'kubernetes_nodes_get', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000010', 'kubernetes_nodes_list', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000011', 'kubernetes_persistentvolumes_get', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000012', 'kubernetes_persistentvolumes_list', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000013', 'kubernetes_secrets_get', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000014', 'kubernetes_secrets_list', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000015', 'kubernetes_endpoints_get', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000016', 'kubernetes_endpoints_list', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000017', 'kubernetes_persistentvolumeclaims_get', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000018', 'kubernetes_persistentvolumeclaims_list', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000019', 'kubernetes_services_get', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000020', 'kubernetes_services_list', NULL, 'v1', 'predefined');

-- k8s.group_version: apps/v1
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000001001', 'kubernetes_deployments_get', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000001002', 'kubernetes_deployments_list', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000001003', 'kubernetes_statefulsets_get', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000001004', 'kubernetes_statefulsets_list', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000001005', 'kubernetes_daemonsets_get', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000001006', 'kubernetes_daemonsets_list', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000001007', 'kubernetes_replicasets_get', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000001008', 'kubernetes_replicasets_list', NULL, 'v1', 'predefined');

-- k8s.group_version: networking.k8s.io/v1
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000002001', 'kubernetes_ingresses_get', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '00000000000000000000000000002002', 'kubernetes_ingresses_list', NULL, 'v1', 'predefined');

-- p8s
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '10000000000000000000000000000001', 'prometheus_query', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '10000000000000000000000000000002', 'prometheus_query_range', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '10000000000000000000000000000003', 'prometheus_alerts', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '10000000000000000000000000000004', 'prometheus_rules', NULL, 'v1', 'predefined');

-- helm
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '20000000000000000000000000000001', 'helm_install', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '20000000000000000000000000000002', 'helm_uninstall', NULL, 'v1', 'predefined');
INSERT INTO `template` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `origin`) VALUES (NULL, NULL, NULL, '20000000000000000000000000000003', 'helm_upgrade', NULL, 'v1', 'predefined');

-- template_command
-- k8s.group_version: v1
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000001', 'kubernetes_pods_get_0', NULL, 'v1', '00000000000000000000000000000001', 0, 'kubernetes.pods.get.v1', '{"namespace":"", "name":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000002', 'kubernetes_pods_list_0', NULL, 'v1', '00000000000000000000000000000002', 0, 'kubernetes.pods.list.v1', '{"namespace":"", "labels":{}}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000003', 'kubernetes_namespaces_get_0', NULL, 'v1', '00000000000000000000000000000003', 0, 'kubernetes.namespaces.get.v1', '{"name":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000004', 'kubernetes_namespaces_list_0', NULL, 'v1', '00000000000000000000000000000004', 0, 'kubernetes.namespaces.list.v1', '{"labels":{}}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000005', 'kubernetes_configmaps_get_0', NULL, 'v1', '00000000000000000000000000000005', 0, 'kubernetes.configmaps.get.v1', '{"namespace":"", "name":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000006', 'kubernetes_configmaps_list_0', NULL, 'v1', '00000000000000000000000000000006', 0, 'kubernetes.configmaps.list.v1', '{"namespace":"", "labels":{}}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000007', 'kubernetes_events_get_0', NULL, 'v1', '00000000000000000000000000000007', 0, 'kubernetes.events.get.v1', '{"namespace":"", "name":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000008', 'kubernetes_events_list_0', NULL, 'v1', '00000000000000000000000000000008', 0, 'kubernetes.events.list.v1', '{"namespace":"", "labels":{}}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000009', 'kubernetes_nodes_get_0', NULL, 'v1', '00000000000000000000000000000009', 0, 'kubernetes.nodes.get.v1', '{"name":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000010', 'kubernetes_nodes_list_0', NULL, 'v1', '00000000000000000000000000000010', 0, 'kubernetes.nodes.list.v1', '{"labels":{}}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000011', 'kubernetes_persistentvolumes_get_0', NULL, 'v1', '00000000000000000000000000000011', 0, 'kubernetes.persistentvolumes.get.v1', '{"name":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000012', 'kubernetes_persistentvolumes_list_0', NULL, 'v1', '00000000000000000000000000000012', 0, 'kubernetes.persistentvolumes.list.v1', '{"labels":{}}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000013', 'kubernetes_secrets_get_0', NULL, 'v1', '00000000000000000000000000000013', 0, 'kubernetes.secrets.get.v1', '{"namespace": "","name":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000014', 'kubernetes_secrets_list_0', NULL, 'v1', '00000000000000000000000000000014', 0, 'kubernetes.secrets.list.v1', '{"namespace":"","labels":{}}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000015', 'kubernetes_endpoints_get_0', NULL, 'v1', '00000000000000000000000000000015', 0, 'kubernetes.endpoints.get.v1', '{"namespace": "","name":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000016', 'kubernetes_endpoints_list_0', NULL, 'v1', '00000000000000000000000000000016', 0, 'kubernetes.endpoints.list.v1', '{"namespace":"","labels":{}}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000017', 'kubernetes_persistentvolumeclaims_get_0', NULL, 'v1', '00000000000000000000000000000017', 0, 'kubernetes.persistentvolumeclaims.get.v1', '{"namespace": "","name":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000018', 'kubernetes_persistentvolumeclaims_list_0', NULL, 'v1', '00000000000000000000000000000018', 0, 'kubernetes.persistentvolumeclaims.list.v1', '{"namespace":"","labels":{}}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000019', 'kubernetes_services_get_0', NULL, 'v1', '00000000000000000000000000000019', 0, 'kubernetes.services.get.v1', '{"namespace": "","name":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000000020', 'kubernetes_services_list_0', NULL, 'v1', '00000000000000000000000000000020', 0, 'kubernetes.services.list.v1', '{"namespace":"","labels":{}}', NULL);

-- k8s.group_version: apps/v1
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000001001', 'kubernetes_deployments_get_0', NULL, 'v1', '00000000000000000000000000001001', 0, 'kubernetes.deployments.get.apps/v1', '{"namespace":"","name":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000001002', 'kubernetes_deployments_list_0', NULL, 'v1', '00000000000000000000000000001002', 0, 'kubernetes.deployments.list.apps/v1', '{"namespace":"","labels":{}}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000001003', 'kubernetes_statefulsets_get_0', NULL, 'v1', '00000000000000000000000000001003', 0, 'kubernetes.statefulsets.get.apps/v1', '{"namespace": "","name":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000001004', 'kubernetes_statefulsets_list_0', NULL, 'v1', '00000000000000000000000000001004', 0, 'kubernetes.statefulsets.list.apps/v1', '{"namespace":"","labels":{}}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000001005', 'kubernetes_daemonsets_get_0', NULL, 'v1', '00000000000000000000000000001005', 0, 'kubernetes.daemonsets.get.apps/v1', '{"namespace": "","name":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000001006', 'kubernetes_daemonsets_list_0', NULL, 'v1', '00000000000000000000000000001006', 0, 'kubernetes.daemonsets.list.apps/v1', '{"namespace":"","labels":{}}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000001007', 'kubernetes_replicasets_get_0', NULL, 'v1', '00000000000000000000000000001007', 0, 'kubernetes.replicasets.get.apps/v1', '{"namespace": "","name":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000001008', 'kubernetes_replicasets_list_0', NULL, 'v1', '00000000000000000000000000001008', 0, 'kubernetes.replicasets.list.apps/v1', '{"namespace":"","labels":{}}', NULL);

-- k8s.group_version: networking.k8s.io/v1
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000002001', 'kubernetes_ingresses_get_0', NULL, 'v1', '00000000000000000000000000002001', 0, 'kubernetes.ingresses.get.networking.k8s.io/v1', '{"namespace": "","name":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '00000000000000000000000000002002', 'kubernetes_ingresses_list_0', NULL, 'v1', '00000000000000000000000000002002', 0, 'kubernetes.ingresses.list.networking.k8s.io/v1', '{"namespace":"","labels":{}}', NULL);

-- p8s
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '10000000000000000000000000000001', 'prometheus_query_0', NULL, 'v1', '10000000000000000000000000000001', 0, 'prometheus.query.v1', '{"url":"","query":"","time":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '10000000000000000000000000000002', 'prometheus_query_range_0', NULL, 'v1', '10000000000000000000000000000002', 0, 'prometheus.query_range.v1', '{"url":"","query":"","start":"","end":"","step":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '10000000000000000000000000000003', 'prometheus_alerts_0', NULL, 'v1', '10000000000000000000000000000003', 0, 'prometheus.alerts.v1', '{"url":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '10000000000000000000000000000004', 'prometheus_rules_0', NULL, 'v1', '10000000000000000000000000000004', 0, 'prometheus.rules.v1', '{"url":""}', NULL);

-- helm
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '20000000000000000000000000000001', 'helm_install_0', NULL, 'v1', '20000000000000000000000000000001', 0, 'helm.install', '{"name":"","chart_name":"","repo_url":"","namespace":"","chart_version":"","values":{}}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '20000000000000000000000000000002', 'helm_uninstall_0', NULL, 'v1', '20000000000000000000000000000002', 0, 'helm.uninstall', '{"name":"","namespace":""}', NULL);
INSERT INTO `template_command` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `sequence`, `method`, `args`, `result_filter`) VALUES (NULL, NULL, NULL, '20000000000000000000000000000003', 'helm_upgrade_0', NULL, 'v1', '20000000000000000000000000000003', 0, 'helm.upgrade', '{"name":"","chart_name":"","repo_url":"","namespace":"","chart_version":"","values":{}}', NULL);
