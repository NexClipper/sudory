
--
-- Table structure for table `service_result_v2`
--

DROP TABLE IF EXISTS `service_result_v2`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `service_result_v2` (
  `pdate` date NOT NULL,
  `cluster_uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `created` datetime(6) NOT NULL,
  `result_type` int(11) NOT NULL,
  `result` longtext COLLATE utf8mb4_bin DEFAULT NULL,
  PRIMARY KEY (`pdate`,`cluster_uuid`,`uuid`),
  KEY `cluster_uuid__uuid` (`cluster_uuid`,`uuid`),
  KEY `cluster_uuid__created` (`cluster_uuid`,`created`),
  KEY `cluster_uuid` (`cluster_uuid`),
  KEY `uuid` (`uuid`),
  KEY `created` (`created`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
 PARTITION BY HASH (to_days(`pdate`))
PARTITIONS 5;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `service_status_v2`
--

DROP TABLE IF EXISTS `service_status_v2`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `service_status_v2` (
  `pdate` date NOT NULL,
  `cluster_uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `created` datetime(6) NOT NULL,
	`step_max` INT(10) UNSIGNED NOT NULL,
  `step_seq` INT(10) UNSIGNED NOT NULL,
  `status` tinyint(3) unsigned NOT NULL,
  `started` datetime DEFAULT NULL,
  `ended` datetime DEFAULT NULL,
  `message` text COLLATE utf8mb4_bin DEFAULT NULL,
  PRIMARY KEY (`pdate`,`cluster_uuid`,`uuid`,`created`),
  KEY `cluster_uuid__uuid` (`cluster_uuid`,`uuid`),
  KEY `cluster_uuid__created` (`cluster_uuid`,`created`),
  KEY `cluster_uuid` (`cluster_uuid`),
  KEY `uuid` (`uuid`),
  KEY `created` (`created`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
 PARTITION BY HASH (to_days(`pdate`))
PARTITIONS 5;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `service_v2`
--

DROP TABLE IF EXISTS `service_v2`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `service_v2` (
  `pdate` date NOT NULL,
  `cluster_uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `created` datetime(6) NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `summary` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `template_uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `flow` TEXT COLLATE utf8mb4_bin NOT NULL,
  `inputs` TEXT COLLATE utf8mb4_bin NOT NULL,
	`step_max` INT(10) UNSIGNED NOT NULL,
  `priority` TINYINT(3) UNSIGNED NOT NULL DEFAULT '0',
  `subscribed_channel` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  PRIMARY KEY (`pdate`,`cluster_uuid`,`uuid`),
  KEY `cluster_uuid__uuid` (`cluster_uuid`,`uuid`),
  KEY `cluster_uuid__created` (`cluster_uuid`,`created`),
  KEY `cluster_uuid` (`cluster_uuid`),
  KEY `uuid` (`uuid`),
  KEY `created` (`created`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
 PARTITION BY HASH (to_days(`pdate`))
PARTITIONS 5;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `template_command_v2`
--

DROP TABLE IF EXISTS `template_command_v2`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `template_command_v2` (
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `summary` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `inputs` text COLLATE utf8mb4_bin DEFAULT NULL,
  `outputs` text COLLATE utf8mb4_bin DEFAULT NULL,
	`client_version` INT(11) NOT NULL DEFAULT '0',
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL,
  `deleted` datetime DEFAULT NULL,
  PRIMARY KEY (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `template_v2`
--

DROP TABLE IF EXISTS `template_v2`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `template_v2` (
  `uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `summary` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `flow` text COLLATE utf8mb4_bin NOT NULL,
  `inputs` text COLLATE utf8mb4_bin DEFAULT NULL,
  `origin` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL,
  `deleted` datetime DEFAULT NULL,
  PRIMARY KEY (`uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
/*!40101 SET character_set_client = @saved_cs_client */;
