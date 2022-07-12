
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `managed_channel_notifier_slackhook` (
  `uuid` char(32) NOT NULL COMMENT 'uuid',
  `url` varchar(255) NOT NULL,
  `request_timeout` int(10) unsigned DEFAULT NULL COMMENT 'request_timeout * time.Second',
  `created` datetime NOT NULL COMMENT 'created',
  `updated` datetime DEFAULT NULL COMMENT 'updated',
  PRIMARY KEY (`uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

