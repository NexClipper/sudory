
--
-- Table structure for table `managed_channel`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `managed_channel` (
  `uuid` char(32) NOT NULL COMMENT 'uuid',
  `name` varchar(255) NOT NULL COMMENT 'name',
  `summary` varchar(255) DEFAULT NULL COMMENT 'summary',
  `event_category` int(11) NOT NULL COMMENT 'event_category',
  `created` datetime NOT NULL COMMENT 'created',
  `updated` datetime DEFAULT NULL COMMENT 'updated',
  `deleted` datetime DEFAULT NULL COMMENT 'deleted',
  PRIMARY KEY (`uuid`) USING BTREE,
  KEY `event_category` (`event_category`) USING BTREE,
  KEY `deleted` (`deleted`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `managed_channel_format`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `managed_channel_format` (
  `uuid` char(32) NOT NULL COMMENT 'uuid',
  `format_type` int(11) NOT NULL,
  `format_data` text NOT NULL,
  `created` datetime NOT NULL COMMENT 'created',
  `updated` datetime DEFAULT NULL COMMENT 'updated',
  PRIMARY KEY (`uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `managed_channel_notifier_console`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `managed_channel_notifier_console` (
  `uuid` char(32) NOT NULL COMMENT 'uuid',
  `created` datetime NOT NULL COMMENT 'created',
  `updated` datetime DEFAULT NULL COMMENT 'updated',
  PRIMARY KEY (`uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `managed_channel_notifier_edge`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `managed_channel_notifier_edge` (
  `uuid` char(32) NOT NULL COMMENT 'uuid',
  `notifier_type` int(11) DEFAULT NULL COMMENT 'notifier_type',
  `created` datetime NOT NULL COMMENT 'created',
  `updated` datetime DEFAULT NULL COMMENT 'updated',
  PRIMARY KEY (`uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `managed_channel_notifier_rabbitmq`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `managed_channel_notifier_rabbitmq` (
  `uuid` char(32) NOT NULL COMMENT 'uuid',
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
  `created` datetime NOT NULL COMMENT 'created',
  `updated` datetime DEFAULT NULL COMMENT 'updated',
  PRIMARY KEY (`uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `managed_channel_notifier_webhook`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `managed_channel_notifier_webhook` (
  `uuid` char(32) NOT NULL COMMENT 'uuid',
  `method` varchar(255) NOT NULL,
  `url` varchar(255) NOT NULL,
  `request_headers` text DEFAULT NULL,
  `request_timeout` int(10) unsigned DEFAULT NULL COMMENT 'request_timeout * time.Second',
  `created` datetime NOT NULL COMMENT 'created',
  `updated` datetime DEFAULT NULL COMMENT 'updated',
  PRIMARY KEY (`uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `managed_channel_status`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `managed_channel_status` (
  `uuid` char(32) NOT NULL COMMENT 'uuid',
  `created` datetime(6) NOT NULL COMMENT 'created',
  `message` text NOT NULL COMMENT 'message',
  PRIMARY KEY (`uuid`,`created`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `managed_channel_status_option`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `managed_channel_status_option` (
  `uuid` char(32) NOT NULL COMMENT 'uuid',
  `status_max_count` int(10) unsigned NOT NULL COMMENT 'status_max_count',
  `created` datetime NOT NULL COMMENT 'created',
  `updated` datetime DEFAULT NULL COMMENT 'updated',
  PRIMARY KEY (`uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
