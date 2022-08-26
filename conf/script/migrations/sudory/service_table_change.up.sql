
-- move old tables
RENAME TABLE `service` TO `old_service`;
RENAME TABLE `service_status` TO `old_service_status`;
RENAME TABLE `service_step` TO `old_service_step`;
RENAME TABLE `service_step_status` TO `old_service_step_status`;
RENAME TABLE `service_result` TO `old_service_result`;

-- create new tables

--
-- Table structure for table `service`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `service` (
  `cluster_uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `pdate` date NOT NULL,
  `timestamp` datetime(6) NOT NULL,
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
  PRIMARY KEY (`cluster_uuid`,`uuid`,`pdate`,`timestamp`),
  INDEX `service_created` (`cluster_uuid`,`created`)
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
  `cluster_uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `pdate` date NOT NULL,
  `timestamp` datetime(6) NOT NULL,
  `result_type` int(11) NOT NULL,
  `result` longtext COLLATE utf8mb4_bin DEFAULT NULL, 
  PRIMARY KEY (`cluster_uuid`,`uuid`,`pdate`,`timestamp`)
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
  `cluster_uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `uuid` char(32) COLLATE utf8mb4_bin NOT NULL,
  `seq` TINYINT(3) UNSIGNED  NOT NULL,
  `pdate` date NOT NULL,
  `timestamp` datetime(6) NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `summary` varchar(255) COLLATE utf8mb4_bin DEFAULT NULL,
  `method` varchar(255) COLLATE utf8mb4_bin NOT NULL,
  `args` text COLLATE utf8mb4_bin DEFAULT NULL,
  `result_filter` varchar(4096) COLLATE utf8mb4_bin DEFAULT NULL,
  `status` TINYINT(3) UNSIGNED  NOT NULL,
  `started` datetime DEFAULT NULL,
  `ended` datetime DEFAULT NULL, 
  `created` datetime NOT NULL,
  PRIMARY KEY (`cluster_uuid`,`uuid`,`seq`,`pdate`,`timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
 PARTITION BY HASH (to_days(`pdate`))
PARTITIONS 5;
/*!40101 SET character_set_client = @saved_cs_client */;


