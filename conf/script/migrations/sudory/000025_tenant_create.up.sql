-- tenant, tenant_channels, tenant_clusters
--
-- Table structure for table `tenant`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `tenant` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `hash` char(40) NOT NULL,
  `pattern` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `summary` varchar(255) DEFAULT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL,
  `deleted` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `hash` (`hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tenant_channels`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `tenant_channels` (
  `channel_uuid` char(32) NOT NULL,
  `tenant_id` bigint(20) NOT NULL,
  PRIMARY KEY (`channel_uuid`) USING BTREE,
	INDEX `reverse` (`tenant_id`, `channel_uuid`) USING BTREE,
  CONSTRAINT `FK_tenant_channels__channel_uuid` FOREIGN KEY (`channel_uuid`) REFERENCES `managed_channel` (`uuid`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_tenant_channels__tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tenant_clusters`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `tenant_clusters` (
  `cluster_id` bigint(20) NOT NULL,
  `tenant_id` bigint(20) NOT NULL,
  PRIMARY KEY (`cluster_id`) USING BTREE,
	INDEX `reverse` (`tenant_id`, `cluster_id`) USING BTREE,
  CONSTRAINT `FK_tenant_clusters__cluster_id` FOREIGN KEY (`cluster_id`) REFERENCES `cluster` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_tenant_clusters__tenant_id` FOREIGN KEY (`tenant_id`) REFERENCES `tenant` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;


-- default value for tenant
INSERT INTO tenant (`hash`, `pattern`, `name`, `summary`, `created`)
SELECT 'da39a3ee5e6b4b0d3255bfef95601890afd80709', '', 'tenant.default', 'hash=''da39a3ee5e6b4b0d3255bfef95601890afd80709'' pattern=''''', NOW()
ON DUPLICATE KEY UPDATE updated=NOW()
;

-- default value for tenant_clusters
INSERT INTO tenant_clusters (`tenant_id`, `cluster_id`)
SELECT tenant.`id`, cluster.`id` FROM cluster, tenant WHERE tenant.`hash` = 'da39a3ee5e6b4b0d3255bfef95601890afd80709'
ON DUPLICATE KEY UPDATE `tenant_id`=VALUES(`tenant_id`), `cluster_id`=VALUES(`cluster_id`) 
;

-- default value for tenant_channels
INSERT INTO tenant_channels (`tenant_id`, `channel_uuid`)
SELECT tenant.`id`, `uuid` FROM managed_channel, tenant WHERE tenant.`hash` = 'da39a3ee5e6b4b0d3255bfef95601890afd80709'
ON DUPLICATE KEY UPDATE `tenant_id`=VALUES(`tenant_id`), `channel_uuid`=VALUES(`channel_uuid`) 
;