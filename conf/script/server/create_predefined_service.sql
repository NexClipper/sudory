use sudory;

-- k8s services for test
SET @cluster_uuid = '00000000000000000000000000000003';
-- -- v1
-- ---- pods get
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-pods-get-4', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-pods-get-4', NULL, 'v1', @uuid, 0, 'kubernetes.pods.get.v1', '{"namespace":"prom","name":"prometheus-prometheus-kube-prometheus-prometheus-0"}', 1, '', NULL, NULL);

-- ---- pods list
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-pods-list-5', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-pods-list-5', NULL, 'v1', @uuid, 0, 'kubernetes.pods.list.v1', '{"release":"prometheus"}', 1, '', NULL, NULL);

-- ---- namespaces get
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-namespaces-get-6', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-namespaces-get-6', NULL, 'v1', @uuid, 0, 'kubernetes.namespaces.get.v1', '{"name":"prom"}', 1, '', NULL, NULL);

-- ---- namespaces list
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-namespaces-list-7', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-namespaces-list-7', NULL, 'v1', @uuid, 0, 'kubernetes.namespaces.list.v1', '{}', 1, '', NULL, NULL);

-- ---- configmaps get
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-configmaps-get-8', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-configmaps-get-8', NULL, 'v1', @uuid, 0, 'kubernetes.configmaps.get.v1', '{"namespace":"prom","name":"prometheus-kube-prometheus-kubelet"}', 1, '', NULL, NULL);

-- ---- configmaps list
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-configmaps-list-9', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-configmaps-list-9', NULL, 'v1', @uuid, 0, 'kubernetes.configmaps.list.v1', '{}', 1, '', NULL, NULL);

-- ---- events get
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-events-get-10', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-events-get-10', NULL, 'v1', @uuid, 0, 'kubernetes.events.get.v1', '{"namespace":"prom","name":"event1"}', 1, '', NULL, NULL);

-- ----  events list
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-events-list-11', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-events-list-11', NULL, 'v1', @uuid, 0, 'kubernetes.events.list.v1', '{}', 1, '', NULL, NULL);

-- ---- nodes get
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-nodes-get-12', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-nodes-get-12', NULL, 'v1', @uuid, 0, 'kubernetes.nodes.get.v1', '{"name":"wslkindmultinodes-control-plane"}', 1, '', NULL, NULL);

-- ---- nodes list
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-nodes-list-13', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-nodes-list-13', NULL, 'v1', @uuid, 0, 'kubernetes.nodes.list.v1', '{}', 1, '', NULL, NULL);

-- ---- persistentvolumes get
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-persistentvolumes-get-14', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-persistentvolumes-get-14', NULL, 'v1', @uuid, 0, 'kubernetes.persistentvolumes.get.v1', '{"namespace":"prom","name":"pvc-8d250848-31d2-47f1-ba1e-133a0e73d030"}', 1, '', NULL, NULL);

-- ---- persistentvolumes list
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-persistentvolumes-list-15', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-persistentvolumes-list-15', NULL, 'v1', @uuid, 0, 'kubernetes.persistentvolumes.list.v1', '{}', 1, '', NULL, NULL);

-- ---- secrets get
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-secrets-get-16', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-secrets-get-16', NULL, 'v1', @uuid, 0, 'kubernetes.secrets.get.v1', '{"namespace":"prom","name":"prometheus-grafana"}', 1, '', NULL, NULL);

-- ---- secrets list
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-secrets-list-17', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-secrets-list-17', NULL, 'v1', @uuid, 0, 'kubernetes.secrets.list.v1', '{}', 1, '', NULL, NULL);

-- -- apps/v1
-- ---- deployments get
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-deployments-get-18', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-deployments-get-18', NULL, 'v1', @uuid, 0, 'kubernetes.deployments.get.apps/v1', '{"namespace":"prom","name":"prometheus-grafana"}', 1, '', NULL, NULL);

-- ---- deployments list
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-deployments-list-19', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-deployments-list-19', NULL, 'v1', @uuid, 0, 'kubernetes.deployments.list.apps/v1', '{}', 1, '', NULL, NULL);

-- ---- statefulsets get
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-statefulsets-get-20', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-statefulsets-get-20', NULL, 'v1', @uuid, 0, 'kubernetes.statefulsets.get.apps/v1', '{"namespace":"prom","name":"prometheus-prometheus-kube-prometheus-prometheus"}', 1, '', NULL, NULL);

-- ---- statefulsets list
SET @uuid = SYS_GUID();
INSERT INTO `service` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `cluster_uuid`, `step_count`, `step_position`, `type`, `epoch`, `interval`, `status`, `result`) VALUES (NULL, NULL, NULL, @uuid, 'service-statefulsets-list-21', NULL, 'v1', @cluster_uuid, 0, 0, 0, 0, 0, 1, NULL);
INSERT INTO `service_step` (`created`, `updated`, `deleted`, `uuid`, `name`, `summary`, `api_version`, `service_uuid`, `sequence`, `method`, `args`, `status`, `result`, `started`, `ended`) VALUES (NULL, NULL, NULL, SYS_GUID(), 'step-statefulsets-list-21', NULL, 'v1', @uuid, 0, 'kubernetes.statefulsets.list.apps/v1', '{}', 1, '', NULL, NULL);
