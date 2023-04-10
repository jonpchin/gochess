-- MySQL dump 10.13  Distrib 5.7.9, for Win64 (x86_64)
--
-- Host: localhost    Database: gochess
-- ------------------------------------------------------
-- Server version	5.7.9-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `gochess`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `gochess` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `gochess`;


--
-- Table structure for table `forums`
--

DROP TABLE IF EXISTS `forums`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `forums` (
  `id` int(2) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `description` varchar(255) NOT NULL,
  `totalthreads` int(10) NOT NULL,
  `totalposts` int(10) NOT NULL,
  `recentuser` varchar(12) NOT NULL,
  `date` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `games`
--

DROP TABLE IF EXISTS `games`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `games` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `white` varchar(12) NOT NULL,
  `black` varchar(12) NOT NULL,
  `gametype` varchar(15) DEFAULT NULL,
  `rated` varchar(3) NOT NULL,
  `whiterating` smallint(4) NOT NULL,
  `blackrating` smallint(4) NOT NULL,
  `timecontrol` smallint(4) NOT NULL,
  `moves` text NOT NULL,
  `totalmoves` int(3) NOT NULL,
  `result` smallint(1) NOT NULL,
  `status` varchar(20) NOT NULL,
  `date` date NOT NULL,
  `time` time NOT NULL,
  `countrywhite` varchar(15) NOT NULL,
  `countryblack` varchar(15) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=696 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `posts`
--

DROP TABLE IF EXISTS `posts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `posts` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `threadId` int(10) NOT NULL,
  `orderId` int(10) DEFAULT NULL,
  `username` varchar(12) DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `body` text NOT NULL,
  `date` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `postThreadsIndex` (`threadId`,`date`)
) ENGINE=InnoDB AUTO_INCREMENT=47 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `rating`
--

DROP TABLE IF EXISTS `rating`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rating` (
  `username` varchar(12) NOT NULL,
  `bullet` smallint(4) NOT NULL,
  `blitz` smallint(4) NOT NULL,
  `standard` smallint(4) NOT NULL,
  `correspondence` smallint(4) NOT NULL,
  `bulletRD` decimal(7,4) NOT NULL,
  `blitzRD` decimal(7,4) NOT NULL,
  `standardRD` decimal(7,4) NOT NULL,
  `correspondenceRD` decimal(7,4) NOT NULL,
  PRIMARY KEY (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ratinghistory`
--

DROP TABLE IF EXISTS `ratinghistory`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ratinghistory` (
  `username` varchar(12) NOT NULL,
  `bullet` text,
  `blitz` text,
  `standard` text,
  `correspondence` text,
  PRIMARY KEY (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `saved`
--

DROP TABLE IF EXISTS `saved`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `saved` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `white` varchar(12) NOT NULL,
  `black` varchar(12) NOT NULL,
  `gametype` varchar(15) DEFAULT NULL,
  `rated` varchar(3) NOT NULL,
  `whiterating` smallint(4) NOT NULL,
  `blackrating` smallint(4) NOT NULL,
  `blackminutes` smallint(4) NOT NULL,
  `blackseconds` smallint(4) NOT NULL,
  `whiteminutes` smallint(4) NOT NULL,
  `whiteseconds` smallint(4) NOT NULL,
  `timecontrol` smallint(4) NOT NULL,
  `moves` varchar(4000) NOT NULL,
  `totalmoves` int(3) NOT NULL,
  `status` varchar(20) NOT NULL,
  `date` date NOT NULL,
  `time` time NOT NULL,
  `countrywhite` varchar(15) DEFAULT NULL,
  `countryblack` varchar(15) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=74 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `threads`
--

DROP TABLE IF EXISTS `threads`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `threads` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `forumId` int(2) NOT NULL,
  `username` varchar(12) NOT NULL,
  `title` varchar(255) NOT NULL,
  `views` int(10) NOT NULL DEFAULT '0',
  `replies` int(10) NOT NULL DEFAULT '0',
  `lastpost` varchar(12) DEFAULT NULL,
  `locked` varchar(3) NOT NULL DEFAULT 'No',
  `date` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `userinfo`
--

DROP TABLE IF EXISTS `userinfo`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `userinfo` (
  `username` varchar(12) NOT NULL,
  `password` char(64) NOT NULL,
  `date` date NOT NULL,
  `time` time NOT NULL,
  `host` varchar(30) DEFAULT NULL,
  `country` varchar(15) DEFAULT NULL,
  `lastpost` datetime DEFAULT NULL,
  `role` varchar(5) NOT NULL DEFAULT 'user',
  PRIMARY KEY (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;



--
-- Table structure for table `grandmaster`
--

DROP TABLE IF EXISTS `grandmaster`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `grandmaster` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `event` varchar(70) NOT NULL,
  `site` varchar(40) NOT NULL,
  `date` varchar(30) NOT NULL,
  `round` varchar(30) NOT NULL,
  `white` varchar(40) NOT NULL,
  `black` varchar(40) NOT NULL,
  `result` varchar(10) NOT NULL,
  `whiteELO` varchar(5) NOT NULL,
  `blackELO` varchar(5) NOT NULL,
  `ECO` varchar(5) NOT NULL,
  `moves` varchar(9000) NOT NULL,
  `eventdate` varchar(30) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `grandmasterECO` (`ECO`) USING HASH
) ENGINE=InnoDB AUTO_INCREMENT=1856328 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;


/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2018-08-17  0:27:23
