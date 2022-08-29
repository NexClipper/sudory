
DROP TABLE IF EXISTS `cluster_information`;

--
-- Table structure for table `cluster_information`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cluster_information` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `cluster_uuid` char(32) NOT NULL,
  `polling_count` int(10) DEFAULT NULL,
  `polling_offset` datetime DEFAULT NULL,
  `created` datetime NOT NULL,
  `updated` datetime DEFAULT NULL,
  `deleted` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `cluster_uuid` (`cluster_uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
