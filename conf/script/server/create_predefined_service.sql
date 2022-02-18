use sudory;

-- k8s services for test
-- -- v1
-- ---- pods get
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (25, '1', 'service-pods-get-4', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (16, '1', 'step-pods-get-4', NULL, 'v1', '1', 0, 'kubernetes.pods.get.v1', '{"namespace":"prom","name":"prometheus-prometheus-kube-prometheus-prometheus-0"}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- ---- pods list
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (26, '2', 'service-pods-list-5', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (17, '2', 'step-pods-list-5', NULL, 'v1', '2', 0, 'kubernetes.pods.list.v1', '{"release":"prometheus"}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- ---- namespaces get
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (27, '3', 'service-namespaces-get-6', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (18, '3', 'step-namespaces-get-6', NULL, 'v1', '3', 0, 'kubernetes.namespaces.get.v1', '{"name":"prom"}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- ---- namespaces list
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (28, '4', 'service-namespaces-list-7', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (19, '4', 'step-namespaces-list-7', NULL, 'v1', '4', 0, 'kubernetes.namespaces.list.v1', '{}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- ---- configmaps get
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (29, '5', 'service-configmaps-get-8', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (20, '5', 'step-configmaps-get-8', NULL, 'v1', '5', 0, 'kubernetes.configmaps.get.v1', '{"namespace":"prom","name":"prometheus-kube-prometheus-kubelet"}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- ---- configmaps list
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (30, '6', 'service-configmaps-list-9', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (21, '6', 'step-configmaps-list-9', NULL, 'v1', '6', 0, 'kubernetes.configmaps.list.v1', '{}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- ---- events get
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (31, '7', 'service-events-get-10', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (22, '7', 'step-events-get-10', NULL, 'v1', '7', 0, 'kubernetes.events.get.v1', '{"namespace":"prom","name":"event1"}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- ----  events list
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (32, '8', 'service-events-list-11', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (23, '8', 'step-events-list-11', NULL, 'v1', '8', 0, 'kubernetes.events.list.v1', '{}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- ---- nodes get
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (33, '9', 'service-nodes-get-12', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (24, '9', 'step-nodes-get-12', NULL, 'v1', '9', 0, 'kubernetes.nodes.get.v1', '{"name":"wslkindmultinodes-control-plane"}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- ---- nodes list
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (34, '10', 'service-nodes-list-13', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (25, '10', 'step-nodes-list-13', NULL, 'v1', '10', 0, 'kubernetes.nodes.list.v1', '{}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- ---- persistentvolumes get
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (35, '11', 'service-persistentvolumes-get-14', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (26, '11', 'step-persistentvolumes-get-14', NULL, 'v1', '11', 0, 'kubernetes.persistentvolumes.get.v1', '{"namespace":"prom","name":"pvc-8d250848-31d2-47f1-ba1e-133a0e73d030"}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- ---- persistentvolumes list
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (36, '12', 'service-persistentvolumes-list-15', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (27, '12', 'step-persistentvolumes-list-15', NULL, 'v1', '12', 0, 'kubernetes.persistentvolumes.list.v1', '{}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- ---- secrets get
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (37, '13', 'service-secrets-get-16', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (28, '13', 'step-secrets-get-16', NULL, 'v1', '13', 0, 'kubernetes.secrets.get.v1', '{"namespace":"prom","name":"prometheus-grafana"}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- ---- secrets list
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (38, '14', 'service-secrets-list-17', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (29, '14', 'step-secrets-list-17', NULL, 'v1', '14', 0, 'kubernetes.secrets.list.v1', '{}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- -- apps/v1
-- ---- deployments get
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (39, '15', 'service-deployments-get-18', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (30, '15', 'step-deployments-get-18', NULL, 'v1', '15', 0, 'kubernetes.deployments.get.apps/v1', '{"namespace":"prom","name":"prometheus-grafana"}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- ---- deployments list
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (40, '16', 'service-deployments-list-19', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (31, '16', 'step-deployments-list-19', NULL, 'v1', '16', 0, 'kubernetes.deployments.list.apps/v1', '{}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- ---- statefulsets get
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (41, '17', 'service-statefulsets-get-20', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (32, '17', 'step-statefulsets-get-20', NULL, 'v1', '17', 0, 'kubernetes.statefulsets.get.apps/v1', '{"namespace":"prom","name":"prometheus-prometheus-kube-prometheus-prometheus"}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- ---- statefulsets list
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (42, '18', 'service-statefulsets-list-21', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (33, '18', 'step-statefulsets-list-21', NULL, 'v1', '18', 0, 'kubernetes.statefulsets.list.apps/v1', '{}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- p8s services for test
-- -- v1
-- ---- query
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (43, '19', 'service-p8s-query-1', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (34, '19', 'step-p8s-query-1', NULL, 'v1', '19', 0, 'prometheus.query.v1', '{"url":"http://localhost:9090","query":"up","time":"2022-02-17T01:59:51.781Z"}', 1, '', NULL, NULL, NULL, NULL, NULL);

-- ---- query_range
INSERT INTO `service` (`id`, `uuid`, `name`, `summary`, `api_version`, `template_uuid`, `origin_kind`, `origin_uuid`, `cluster_uuid`, `assigned_client_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`, `created`, `updated`, `deleted`) VALUES (44, '20', 'service-p8s-query_range-2', NULL, 'v1', 'template_uuid', 'template', 'template_uuid', '00000000000000000000000000000003', NULL, 0, 0, 0, 0, 0, 1, NULL, NULL, NULL, NULL);
INSERT INTO `service_step` (`id`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`, `created`, `updated`, `deleted`) VALUES (35, '20', 'step-p8s-query_range-2', NULL, 'v1', '20', 0, 'prometheus.query_range.v1', '{"url":"http://localhost:9090","end":"2022-02-17T02:10:30.781Z","query":"rate(prometheus_tsdb_head_samples_appended_total[5m])","start":"2022-02-17T01:10:30.781Z","step":"1y"}', 1, '', NULL, NULL, NULL, NULL, NULL);
