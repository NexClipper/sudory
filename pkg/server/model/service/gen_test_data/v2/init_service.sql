
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
