-- MySQL dump 10.19  Distrib 10.3.32-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: sudory
-- ------------------------------------------------------
-- Server version       10.3.32-MariaDB-0ubuntu0.20.04.1

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
-- Current Database: `sudory`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `sudory` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

USE `sudory`;

--
-- Table structure for table `client`
--

DROP TABLE IF EXISTS `client`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `client` (
  `ID` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `MACHINE_ID` varchar(45) COLLATE utf8_bin NOT NULL,
  `CLUSTER_ID` bigint(20) unsigned NOT NULL,
  `ACTIVE` tinyint(1) NOT NULL DEFAULT 0,
  `IP` varchar(45) COLLATE utf8_bin NOT NULL,
  `PORT` int(11) NOT NULL,
  `CREATED` timestamp NULL DEFAULT current_timestamp(),
  `UPDATED` timestamp NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`ID`),
  KEY `FK_CLIENT_idx` (`CLUSTER_ID`),
  CONSTRAINT `FK_CLIENT` FOREIGN KEY (`CLUSTER_ID`) REFERENCES `cluster` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='client 관리 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cluster`
--

DROP TABLE IF EXISTS `cluster`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cluster` (
  `ID` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `NAME` varchar(200) COLLATE utf8_bin NOT NULL,
  `CREATED` timestamp NULL DEFAULT current_timestamp(),
  `UPDATED` timestamp NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='cluster 관리 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `service`
--

DROP TABLE IF EXISTS `service`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `service` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `deleted` datetime DEFAULT NULL,
  `uuid` char(32) COLLATE utf8_bin NOT NULL,
  `name` varchar(255) COLLATE utf8_bin NOT NULL,
  `summary` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `api_version` varchar(255) COLLATE utf8_bin NOT NULL,
  `cluster_uuid` char(32) COLLATE utf8_bin NOT NULL COMMENT 'clusters uuid',
  `status` int(11) NOT NULL DEFAULT 0 COMMENT 'status',
  `step_count` int(11) DEFAULT 0,
  `step_position` int(11) DEFAULT 0,
  `type` int(11) DEFAULT 0,
  `epoch` int(11) DEFAULT 0,
  `interval` int(11) DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UQE_service_uuid` (`uuid`),
  KEY `IDX_service_cluster_uuid` (`cluster_uuid`),
  KEY `IDX_service_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `service_step`
--

DROP TABLE IF EXISTS `service_step`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `service_step` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `deleted` datetime DEFAULT NULL,
  `uuid` char(32) COLLATE utf8_bin NOT NULL,
  `name` varchar(255) COLLATE utf8_bin NOT NULL,
  `summary` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `api_version` varchar(255) COLLATE utf8_bin NOT NULL,
  `service_uuid` char(32) COLLATE utf8_bin NOT NULL COMMENT 'services uuid',
  `sequence` int(11) NOT NULL COMMENT 'sequence',
  `method` varchar(255) COLLATE utf8_bin NOT NULL COMMENT 'method',
  `args` text COLLATE utf8_bin DEFAULT NULL COMMENT 'args',
  `status` int(11) NOT NULL DEFAULT 1 COMMENT 'status',
  `result` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT 'result',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UQE_service_step_uuid` (`uuid`),
  KEY `IDX_service_step_service_uuid` (`service_uuid`),
  KEY `status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `template`
--

DROP TABLE IF EXISTS `template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `template` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `deleted` datetime DEFAULT NULL,
  `uuid` char(32) COLLATE utf8_bin NOT NULL,
  `name` varchar(255) COLLATE utf8_bin NOT NULL,
  `summary` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `api_version` varchar(255) COLLATE utf8_bin NOT NULL COMMENT 'api version',
  `origin` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT 'origin',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB AUTO_INCREMENT=46 DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='table template';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `template_command`
--

DROP TABLE IF EXISTS `template_command`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `template_command` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created` datetime DEFAULT NULL,
  `updated` datetime DEFAULT NULL,
  `deleted` datetime DEFAULT NULL,
  `uuid` char(32) COLLATE utf8_bin NOT NULL,
  `name` varchar(255) COLLATE utf8_bin NOT NULL,
  `summary` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `api_version` varchar(255) COLLATE utf8_bin NOT NULL,
  `template_uuid` char(32) COLLATE utf8_bin NOT NULL COMMENT 'templates uuid',
  `sequence` int(11) NOT NULL COMMENT 'sequence',
  `method` varchar(255) COLLATE utf8_bin NOT NULL COMMENT 'method',
  `args` text COLLATE utf8_bin DEFAULT NULL COMMENT 'arguments',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uuid` (`uuid`) USING BTREE,
  KEY `template_uuid` (`template_uuid`)
) ENGINE=InnoDB AUTO_INCREMENT=52 DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='table template command';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `token`
--

DROP TABLE IF EXISTS `token`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `token` (
  `CLUSTER_ID` bigint(20) unsigned NOT NULL,
  `KEY` varchar(100) COLLATE utf8_bin NOT NULL,
  `CREATED` timestamp NULL DEFAULT current_timestamp(),
  `UPDATED` timestamp NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`CLUSTER_ID`),
  CONSTRAINT `FK_TOKEN` FOREIGN KEY (`CLUSTER_ID`) REFERENCES `cluster` (`ID`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='TOKEN 관리 테이블';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-01-18 10:58:56
