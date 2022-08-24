
-- DROP DATABASE IF EXISTS sudory_schema_test_v3;
-- CREATE DATABASE sudory_schema_test_v3;

--
-- Table structure for table `service`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `service` (
  `pdate` date NOT NULL,
  `cluster_uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `revision` TINYINT(3) UNSIGNED NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `summary` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `template_uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `step_count` TINYINT(3) UNSIGNED NOT NULL,
  `subscribed_channel` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `assigned_client_uuid` char(32) COLLATE utf8mb4_bin DEFAULT NULL,
  `step_position` TINYINT(3) UNSIGNED NOT NULL,
  `status` TINYINT(3) UNSIGNED NOT NULL,
  `message` text COLLATE utf8mb4_bin DEFAULT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`pdate`,`cluster_uuid`,`uuid`,`revision`)
  INDEX `cluster_status` (`cluster_uuid`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
 PARTITION BY HASH (to_days(`pdate`))
PARTITIONS 5;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `service_result`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `service_result` (
  `pdate` date NOT NULL,
  `cluster_uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `revision` TINYINT(3) UNSIGNED NOT NULL,
  `result_type` int(11) NOT NULL,
  `result` longtext COLLATE utf8mb4_bin DEFAULT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`pdate`,`cluster_uuid`,`uuid`,`revision`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
 PARTITION BY HASH (to_days(`pdate`))
PARTITIONS 5;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `service_step`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `service_step` (
  `pdate` date NOT NULL,
  `cluster_uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `seq` TINYINT(3) UNSIGNED  NOT NULL,
  `revision` TINYINT(3) UNSIGNED NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `summary` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `method` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `args` text COLLATE utf8mb4_bin DEFAULT NULL,
  `result_filter` varchar(4096) COLLATE utf8mb4_bin DEFAULT NULL,
  `status` TINYINT(3) UNSIGNED  NOT NULL,
  `started` datetime DEFAULT NULL,
  `ended` datetime DEFAULT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL,
  PRIMARY KEY (`pdate`,`cluster_uuid`,`uuid`,`seq`,`revision`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
 PARTITION BY HASH (to_days(`pdate`))
PARTITIONS 5;
/*!40101 SET character_set_client = @saved_cs_client */;
