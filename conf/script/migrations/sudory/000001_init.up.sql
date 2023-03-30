-- MariaDB dump 10.19  Distrib 10.8.3-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: sudory
-- ------------------------------------------------------
-- Server version       10.8.3-MariaDB-1:10.8.3+maria~jammy

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

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
  `created` datetime DEFAULT NULL COMMENT 'created',
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
  `created` datetime DEFAULT NULL COMMENT 'created',
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
  `created` datetime DEFAULT NULL COMMENT 'created',
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
  `created` datetime DEFAULT NULL COMMENT 'created',
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
  `created` datetime DEFAULT NULL COMMENT 'created',
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
  `created` datetime DEFAULT NULL COMMENT 'created',
  `updated` datetime DEFAULT NULL COMMENT 'updated',
  `deleted` datetime DEFAULT NULL COMMENT 'deleted',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `UQE_uuid` (`uuid`) USING BTREE,
  KEY `IDX_deleted` (`deleted`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cluster`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `cluster` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uuid` char(32) COLLATE utf8mb3_bin NOT NULL,
  `name` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `summary` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `polling_option` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `polling_limit` smallint(6) NOT NULL DEFAULT 0,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `deleted` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin COMMENT='cluster 관리 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cluster_token`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `cluster_token` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uuid` char(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL,
  `summary` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin DEFAULT NULL,
  `cluster_uuid` char(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `token` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL COMMENT 'token',
  `issued_at_time` datetime NOT NULL COMMENT 'issued at time',
  `expiration_time` datetime NOT NULL COMMENT 'expiration time',
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `deleted` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uuid` (`uuid`) USING BTREE,
  UNIQUE KEY `token` (`token`),
  KEY `IDX_token_cluster_uuid` (`cluster_uuid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `global_variant`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `global_variant` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uuid` char(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL,
  `summary` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin DEFAULT NULL,
  `value` text CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL COMMENT 'value',
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `deleted` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `service`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `service` (
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
-- Table structure for table `service_result`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `service_result` (
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
-- Table structure for table `service_status`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `service_status` (
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
-- Table structure for table `service_step`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `service_step` (
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
-- Table structure for table `service_step_status`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `service_step_status` (
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

--
-- Table structure for table `session`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `session` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uuid` char(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL,
  `cluster_uuid` char(32) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL COMMENT 'user_uuid',
  `token` text CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL COMMENT 'token',
  `issued_at_time` datetime NOT NULL COMMENT 'issued as time',
  `expiration_time` datetime NOT NULL COMMENT 'expiration time',
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `deleted` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uuid` (`uuid`) USING BTREE,
  KEY `cluster_uuid` (`cluster_uuid`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `template`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `template` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uuid` char(32) COLLATE utf8mb3_bin NOT NULL,
  `name` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `summary` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `origin` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT 'origin',
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `deleted` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin COMMENT='table template';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `template_command`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `template_command` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uuid` char(32) COLLATE utf8mb3_bin NOT NULL,
  `name` varchar(255) COLLATE utf8mb3_bin NOT NULL,
  `summary` varchar(255) COLLATE utf8mb3_bin DEFAULT NULL,
  `template_uuid` char(32) COLLATE utf8mb3_bin NOT NULL COMMENT 'templates uuid',
  `sequence` int(11) NOT NULL COMMENT 'sequence',
  `method` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT 'method',
  `args` text COLLATE utf8mb3_bin DEFAULT NULL COMMENT 'arguments',
  `result_filter` varchar(4096) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'result_filter',
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `deleted` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uuid` (`uuid`) USING BTREE,
  KEY `template_uuid` (`template_uuid`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin COMMENT='table template command';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `template_recipe`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `template_recipe` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `summary` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `method` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `args` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `deleted` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `method_v2` (`method`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
