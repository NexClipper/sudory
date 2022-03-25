--	https://docs.google.com/spreadsheets/d/1vp8PuyQanxLlfmrKMM6BPI8Mf02lZmd9X2YsVEM4A80/edit#gid=1701104947
	use sudory;
--	template_recipe
-- kubernetes	
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.configmaps.get.v1', 'namespace,name', 'kubernetes-configmaps-get-v1', 'namespace:required;string;,name:required;string;', '1', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.configmaps.list.v1', 'namespace,labels', 'kubernetes-configmaps-list-v1', 'namespace:optional;string;,labels:optional;key-value;', '2', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.events.get.v1', 'namespace,name', 'kubernetes-events-get-v1', 'namespace:required;string;,name:required;string;', '3', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.events.list.v1', 'namespace,labels', 'kubernetes-events-list-v1', 'namespace:optional;string;,labels:optional;key-value;', '4', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.namespaces.get.v1', 'name', 'kubernetes-namespaces-get-v1', 'name:required;string;', '5', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.namespaces.list.v1', 'labels', 'kubernetes-namespaces-list-v1', 'labels:optional;key-value;', '6', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.nodes.get.v1', 'name', 'kubernetes-nodes-get-v1', 'name:required;string;', '7', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.nodes.list.v1', 'labels', 'kubernetes-nodes-list-v1', 'labels:optional;key-value;', '8', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.persistentvolumes.get.v1', 'name', 'kubernetes-persistentvolumes-get-v1', 'name:required;string;', '9', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.persistentvolumes.list.v1', 'labels', 'kubernetes-persistentvolumes-list-v1', 'labels:optional;key-value;', '10', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.pods.get.v1', 'namespace,name', 'kubernetes-pods-get-v1', 'namespace:required;string;,name:required;string;', '11', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.pods.list.v1', 'namespace,labels', 'kubernetes-pods-list-v1', 'namespace:optional;string;,labels:optional;key-value;', '12', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.secrets.get.v1', 'namespace,name', 'kubernetes-secrets-get-v1', 'namespace:required;string;,name:required;string;', '13', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.secrets.list.v1', 'namespace,labels', 'kubernetes-secrets-list-v1', 'namespace:optional;string;,labels:optional;key-value;', '14', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.endpoints.get.v1', 'namespace,name', 'kubernetes-endpoints-get-v1', 'namespace:required;string;,name:required;string;', '15', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.endpoints.list.v1', 'namespace,labels', 'kubernetes-endpoints-list-v1', 'namespace:optional;string;,labels:optional;key-value;', '16', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.persistentvolumeclaims.get.v1', 'namespace,name', 'kubernetes-persistentvolumeclaims-get-v1', 'namespace:required;string;,name:required;string;', '17', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.persistentvolumeclaims.list.v1', 'namespace,labels', 'kubernetes-persistentvolumeclaims-list-v1', 'optional;string;', '18', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.services.get.v1', 'namespace,name', 'kubernetes-services-get-v1', 'namespace:required;string;,name:required;string;', '19', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.services.list.v1', 'namespace,labels', 'kubernetes-services-list-v1', 'namespace:optional;string;,labels:optional;key-value;', '20', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.deployments.get.apps/v1', 'namespace,name', 'kubernetes-deployments-get-apps/v1', 'namespace:required;string;,name:required;string;', '21', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.deployments.list.apps/v1', 'namespace,labels', 'kubernetes-deployments-list-apps/v1', 'namespace:optional;string;,labels:optional;key-value;', '22', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.statefulsets.get.apps/v1', 'namespace,name', 'kubernetes-statefulsets-get-apps/v1', 'namespace:required;string;,name:required;string;', '23', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.statefulsets.list.apps/v1', 'namespace,labels', 'kubernetes-statefulsets-list-apps/v1', 'namespace:optional;string;,labels:optional;key-value;', '24', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.daemonsets.get.apps/v1', 'namespace,name', 'kubernetes-daemonsets-get-apps/v1', 'namespace:required;string;,name:required;string;', '25', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.daemonsets.list.apps/v1', 'namespace,labels', 'kubernetes-daemonsets-list-apps/v1', 'namespace:optional;string;,labels:optional;key-value;', '26', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.replicasets.get.apps/v1', 'namespace,name', 'kubernetes-replicasets-get-apps/v1', 'namespace:required;string;,name:required;string;', '27', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.replicasets.list.apps/v1', 'namespace,labels', 'kubernetes-replicasets-list-apps/v1', 'namespace:optional;string;,labels:optional;key-value;', '28', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.ingresses.get.networking.k8s.io/v1', 'namespace,name', 'kubernetes-ingresses-get-networking.k8s.io/v1', 'namespace:required;string;,name:required;string;', '29', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.ingresses.list.networking.k8s.io/v1', 'namespace,labels', 'kubernetes-ingresses-list-networking.k8s.io/v1', 'namespace:optional;string;,labels:optional;key-value;', '30', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.storageclasses.get.storage.k8s.io/v1', 'name', 'kubernetes-storageclasses-get-storage.k8s.io/v1', 'name:required;string;', '31', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.storageclasses.list.storage.k8s.io/v1', 'labels', 'kubernetes-storageclasses-list-storage.k8s.io/v1', 'labels:optional;key-value;', '32', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.prometheuses.get.monitoring.coreos.com/v1', 'namespace,name', 'kubernetes-prometheuses-get-monitoring.coreos.com/v1', 'namespace:required;string;,name:required;string;', '33', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.prometheuses.list.monitoring.coreos.com/v1', 'namespace,labels', 'kubernetes-prometheuses-list-monitoring.coreos.com/v1', 'namespace:optional;string;,labels:optional;key-value;', '34', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.prometheusrules.get.monitoring.coreos.com/v1', 'namespace,name', 'kubernetes-prometheusrules-get-monitoring.coreos.com/v1', 'namespace:required;string;,name:required;string;', '35', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.prometheusrules.list.monitoring.coreos.com/v1', 'namespace,labels', 'kubernetes-prometheusrules-list-monitoring.coreos.com/v1', 'namespace:optional;string;,labels:optional;key-value;', '36', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.servicemonitors.get.monitoring.coreos.com/v1', 'namespace,name', 'kubernetes-servicemonitors-get-monitoring.coreos.com/v1', 'namespace:required;string;,name:required;string;', '37', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('kubernetes.servicemonitors.list.monitoring.coreos.com/v1', 'namespace,labels', 'kubernetes-servicemonitors-list-monitoring.coreos.com/v1', 'namespace:optional;string;,labels:optional;key-value;', '38', '2022-03-28 14:05:07');
-- prometheus	
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('prometheus.query..v1', 'url,query,time', 'prometheus-query--v1', 'url:required;string;,query:required;string;,time:optional;string;', '39', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('prometheus.query_range..v1', 'url,query,start,end,step', 'prometheus-query_range--v1', 'url:required;string;,query:required;string;,start:required;string;,end:required;string;,step:required;string;', '40', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('prometheus.alerts..v1', 'url', 'prometheus-alerts--v1', 'url:required;string;', '41', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('prometheus.rules..v1', 'url', 'prometheus-rules--v1', 'url:required;string;', '42', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('prometheus.alertmanagers..v1', 'url', 'prometheus-alertmanagers--v1', 'url:required;string;', '43', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('prometheus.targets..v1', 'url', 'prometheus-targets--v1', 'url:required;string;', '44', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('prometheus.targets/metadata..v1', 'url,match_target,metric,limit', 'prometheus-targets/metadata--v1', 'url:required;string;,match_target:optional;string;,metric:optional;string;,limit:optional;string;', '45', '2022-03-28 14:05:07');
-- jq	
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('jq...', 'filter,input', 'jq---', 'filter:string;,input:map[string]interface{};', '46', '2022-03-28 14:05:07');
-- helm	
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('helm.install..', 'name,chart_name,repo_url,namespace,chart_version,values', 'helm-install--', 'name:require;string;,chart_name:require;string;,repo_url:require;string;,namespace:require;string;,chart_version:optional;string;,values:optional;key-value;', '47', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('helm.uninstall..', 'name,namespace', 'helm-uninstall--', 'name:require;string;,namespace:require;string;', '48', '2022-03-28 14:05:07');
	REPLACE INTO `template_recipe` ( `method`, `args`, `name`, `summary`, `id`, `created`) VALUES ('helm.upgrade..', 'name,chart_name,repo_url,namespace,chart_version,values', 'helm-upgrade--', 'name:require;string;,chart_name:require;string;,repo_url:require;string;,namespace:require;string;,chart_version:optional;string;,values:optional;key-value;', '49', '2022-03-28 14:05:07');
	