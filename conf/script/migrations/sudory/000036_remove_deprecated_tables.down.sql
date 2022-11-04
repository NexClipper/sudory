
-- DROP TABLES (DOWN)

--
-- Table structure for table `channel`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `channel` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `uuid` char(32) NOT NULL COMMENT 'uuid',
  `name` varchar(255) NOT NULL COMMENT 'name',
  `summary` varchar(255) DEFAULT NULL COMMENT 'summary',
  `cluster_uuid` char(32) NOT NULL,
  `created` datetime NOT NULL COMMENT 'created',
  `updated` datetime DEFAULT NULL COMMENT 'updated',
  `deleted` datetime DEFAULT NULL COMMENT 'deleted',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `UQE_uuid` (`uuid`) USING BTREE,
  KEY `IDX_channel_cluster_uuid` (`cluster_uuid`) USING BTREE,
  KEY `IDX_deleted` (`deleted`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `channel_notifier_console`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `channel_notifier_console` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `uuid` char(32) NOT NULL COMMENT 'uuid',
  `name` varchar(255) NOT NULL COMMENT 'name',
  `summary` varchar(255) DEFAULT NULL COMMENT 'summary',
  `content_type` varchar(255) NOT NULL,
  `created` datetime NOT NULL COMMENT 'created',
  `updated` datetime DEFAULT NULL COMMENT 'updated',
  `deleted` datetime DEFAULT NULL COMMENT 'deleted',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `UQE_uuid` (`uuid`) USING BTREE,
  KEY `IDX_deleted` (`deleted`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `channel_notifier_edge`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `channel_notifier_edge` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `channel_uuid` char(32) NOT NULL,
  `notifier_type` varchar(255) NOT NULL,
  `notifier_uuid` char(32) NOT NULL,
  `created` datetime NOT NULL COMMENT 'created',
  `updated` datetime DEFAULT NULL COMMENT 'updated',
  `deleted` datetime DEFAULT NULL COMMENT 'deleted',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `IDX_channel_notifier_edge_notifier_type` (`notifier_type`) USING BTREE,
  KEY `IDX_channel_notifier_edge_notifier_uuid` (`notifier_uuid`) USING BTREE,
  KEY `IDX_deleted` (`deleted`) USING BTREE,
  KEY `IDX_channel_notifier_edge_channel_uuid` (`channel_uuid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `channel_notifier_rabbitmq`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `channel_notifier_rabbitmq` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `uuid` char(32) NOT NULL COMMENT 'uuid',
  `name` varchar(255) NOT NULL COMMENT 'name',
  `summary` varchar(255) DEFAULT NULL COMMENT 'summary',
  `url` varchar(255) NOT NULL,
  `exchange` varchar(255) DEFAULT NULL,
  `routing_key` varchar(255) DEFAULT NULL,
  `mandatory` tinyint(1) DEFAULT NULL,
  `immediate` tinyint(1) DEFAULT NULL,
  `message_headers` text DEFAULT NULL,
  `message_content_type` varchar(255) DEFAULT NULL,
  `message_content_encoding` varchar(255) DEFAULT NULL,
  `message_delivery_mode` int(10) unsigned DEFAULT NULL,
  `message_priority` int(10) unsigned DEFAULT NULL,
  `message_correlation_id` varchar(255) DEFAULT NULL,
  `message_reply_to` varchar(255) DEFAULT NULL,
  `message_expiration` varchar(255) DEFAULT NULL,
  `message_message_id` varchar(255) DEFAULT NULL,
  `message_timestamp` tinyint(1) DEFAULT NULL,
  `message_type` varchar(255) DEFAULT NULL,
  `message_user_id` varchar(255) DEFAULT NULL,
  `message_app_id` varchar(255) DEFAULT NULL,
  `content_type` varchar(255) NOT NULL,
  `created` datetime NOT NULL COMMENT 'created',
  `updated` datetime DEFAULT NULL COMMENT 'updated',
  `deleted` datetime DEFAULT NULL COMMENT 'deleted',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `UQE_uuid` (`uuid`) USING BTREE,
  KEY `IDX_deleted` (`deleted`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `channel_notifier_status`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `channel_notifier_status` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `uuid` char(32) NOT NULL COMMENT 'uuid',
  `notifier_type` varchar(255) NOT NULL,
  `notifier_uuid` char(32) NOT NULL,
  `error` text DEFAULT NULL,
  `created` datetime NOT NULL COMMENT 'created',
  `updated` datetime DEFAULT NULL COMMENT 'updated',
  `deleted` datetime DEFAULT NULL COMMENT 'deleted',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `UQE_uuid` (`uuid`) USING BTREE,
  KEY `IDX_channel_notifier_status_notifier_uuid` (`notifier_uuid`) USING BTREE,
  KEY `IDX_deleted` (`deleted`) USING BTREE,
  KEY `IDX_channel_notifier_status_notifier_type` (`notifier_type`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `channel_notifier_webhook`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `channel_notifier_webhook` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `uuid` char(32) NOT NULL COMMENT 'uuid',
  `name` varchar(255) NOT NULL COMMENT 'name',
  `summary` varchar(255) DEFAULT NULL COMMENT 'summary',
  `method` varchar(255) NOT NULL,
  `url` varchar(255) NOT NULL,
  `request_headers` text DEFAULT NULL,
  `request_timeout` varchar(16) DEFAULT NULL COMMENT 'fmt(time.ParseDuration)',
  `content_type` varchar(255) NOT NULL,
  `created` datetime NOT NULL COMMENT 'created',
  `updated` datetime DEFAULT NULL COMMENT 'updated',
  `deleted` datetime DEFAULT NULL COMMENT 'deleted',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `UQE_uuid` (`uuid`) USING BTREE,
  KEY `IDX_deleted` (`deleted`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `old_service`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `old_service` (
  `uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `created` datetime(6) NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `summary` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `cluster_uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `template_uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `step_count` int(11) NOT NULL,
  `subscribed_channel` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  PRIMARY KEY (`uuid`,`created`),
  KEY `cluster_uuid` (`cluster_uuid`),
  KEY `template_uuid` (`template_uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
 PARTITION BY HASH (to_days(`created`))
PARTITIONS 5;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `old_service_result`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `old_service_result` (
  `uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `created` datetime(6) NOT NULL,
  `result_type` int(11) NOT NULL,
  `result` longtext COLLATE utf8mb4_bin DEFAULT NULL,
  PRIMARY KEY (`uuid`,`created`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
 PARTITION BY HASH (to_days(`created`))
PARTITIONS 5;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `old_service_status`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `old_service_status` (
  `uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `created` datetime(6) NOT NULL,
  `assigned_client_uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `step_position` int(11) NOT NULL,
  `status` int(11) NOT NULL,
  `message` text COLLATE utf8mb4_bin DEFAULT NULL,
  PRIMARY KEY (`uuid`,`created`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
 PARTITION BY HASH (to_days(`created`))
PARTITIONS 5;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `old_service_step`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `old_service_step` (
  `uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `sequence` int(11) NOT NULL,
  `created` datetime(6) NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `summary` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `method` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `args` text COLLATE utf8mb4_bin DEFAULT NULL,
  `result_filter` varchar(4096) COLLATE utf8mb4_bin DEFAULT NULL,
  PRIMARY KEY (`uuid`,`sequence`,`created`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
 PARTITION BY HASH (to_days(`created`))
PARTITIONS 5;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `old_service_step_status`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `old_service_step_status` (
  `uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `sequence` int(11) NOT NULL,
  `created` datetime(6) NOT NULL,
  `status` int(11) NOT NULL,
  `started` datetime DEFAULT NULL,
  `ended` datetime DEFAULT NULL,
  PRIMARY KEY (`uuid`,`sequence`,`created`),
  KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
 PARTITION BY HASH (to_days(`created`))
PARTITIONS 5;
/*!40101 SET character_set_client = @saved_cs_client */;