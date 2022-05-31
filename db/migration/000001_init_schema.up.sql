-- MySQL dump 10.13  Distrib 8.0.28, for Linux (x86_64)
--
-- Host: localhost    Database: gopool
-- ------------------------------------------------------
-- Server version	8.0.28

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `machine`
--

DROP TABLE IF EXISTS `machine`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `machine` (
  `machine_id` int NOT NULL AUTO_INCREMENT,
  `phone_number` varchar(12) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ip_address` varchar(15) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`machine_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `machine`
--

LOCK TABLES `machine` WRITE;
/*!40000 ALTER TABLE `machine` DISABLE KEYS */;
INSERT INTO `machine` VALUES (1,'744265634','::1');
/*!40000 ALTER TABLE `machine` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pool_config`
--

DROP TABLE IF EXISTS `pool_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pool_config` (
  `pool_id` int NOT NULL,
  `single_vote` tinyint NOT NULL,
  `start_date` timestamp NOT NULL,
  `end_date` timestamp NOT NULL,
  KEY `pool_id` (`pool_id`),
  CONSTRAINT `pool_config_ibfk_1` FOREIGN KEY (`pool_id`) REFERENCES `pools` (`pool_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pool_config`
--

LOCK TABLES `pool_config` WRITE;
/*!40000 ALTER TABLE `pool_config` DISABLE KEYS */;
INSERT INTO `pool_config` VALUES (1,1,'2022-05-12 00:00:00','2022-05-13 00:00:00'),(24,1,'2022-05-06 00:00:00','2022-05-07 00:00:00'),(25,1,'2022-05-13 00:00:00','2022-05-06 00:00:00'),(26,1,'2022-05-26 00:00:00','2022-05-26 00:00:00'),(27,1,'2022-05-20 00:00:00','2022-05-20 00:00:00'),(28,1,'2022-05-22 00:00:00','2022-05-22 00:00:00'),(29,1,'2022-05-23 00:00:00','2022-05-23 00:00:00'),(30,1,'2022-05-23 00:05:00','2022-05-23 00:10:00'),(31,1,'2022-05-23 00:11:00','2022-05-23 00:12:00'),(32,1,'2022-05-23 14:37:00','2022-05-23 15:00:00'),(33,1,'2022-05-23 15:05:00','2022-05-23 16:05:00'),(34,1,'2022-05-23 19:16:00','2022-05-23 19:18:00'),(35,1,'2022-05-25 13:43:00','2022-05-25 13:50:00'),(36,1,'2022-05-25 13:54:00','2022-05-25 14:52:00');
/*!40000 ALTER TABLE `pool_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pool_option`
--

DROP TABLE IF EXISTS `pool_option`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pool_option` (
  `option_id` int NOT NULL AUTO_INCREMENT,
  `pool_id` int NOT NULL,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`option_id`),
  KEY `pool_id` (`pool_id`),
  CONSTRAINT `pool_option_ibfk_1` FOREIGN KEY (`pool_id`) REFERENCES `pools` (`pool_id`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pool_option`
--

LOCK TABLES `pool_option` WRITE;
/*!40000 ALTER TABLE `pool_option` DISABLE KEYS */;
INSERT INTO `pool_option` VALUES (3,1,'name2','desc1'),(6,19,'ddvd','ads'),(7,19,'dvddv sd','devs'),(8,19,'dvsvds','vddd'),(9,20,'sdcsd','dxfs'),(10,21,'ascxs','ascend'),(11,24,'option1','desc1'),(12,25,'wad','dqwqw'),(13,26,'opt1','desc1'),(14,28,'owd1','sdcdw'),(15,29,'psdfc','pdf'),(16,30,'sdc2','dads'),(17,31,'edfwe3','3r2'),(18,32,'test','test desc'),(19,33,'test1','acdc'),(20,34,'13w2','23'),(21,35,'opt1','option1'),(22,35,'opt2','option2'),(23,35,'opt3','option3'),(24,36,'opt1','option1'),(25,36,'opt2','option2'),(26,36,'opt3','option3');
/*!40000 ALTER TABLE `pool_option` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pools`
--

DROP TABLE IF EXISTS `pools`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pools` (
  `pool_id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `nr_of_options` int NOT NULL,
  PRIMARY KEY (`pool_id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `pools_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pools`
--

LOCK TABLES `pools` WRITE;
/*!40000 ALTER TABLE `pools` DISABLE KEYS */;
INSERT INTO `pools` VALUES (1,1,'pool 1_13',1),(2,1,'pool 1',1),(3,1,'pool 2',0),(4,1,'pool 3',0),(5,1,'pool 5',0),(6,1,'pool 5',0),(7,1,'pool 5',0),(8,1,'pool 6',0),(9,1,'pool 6',0),(10,1,'pool7',0),(11,1,'pool8',4),(12,1,'cc',2),(13,1,'ghfhgf',3),(14,1,'ddvd',3),(15,1,'xxcdx',2),(16,1,'dvdsf',2),(17,1,'sdks',3),(18,1,'pool 234',3),(19,1,'fdvdf',3),(20,1,'pool 1',1),(21,1,'pool 1_1',1),(22,1,'pool -2',1),(23,1,'pool test 3',1),(24,1,'pool test 3',1),(25,1,'pool24',1),(26,2,'pool1',1),(27,2,'cvcv',1),(28,3,'pool1',1),(29,3,'pool2',1),(30,3,'pool3',1),(31,3,'pool4',1),(32,3,'test1',1),(33,3,'test2',1),(34,3,'test3',1),(35,3,'test 13:41',3),(36,3,'test 13:52',3);
/*!40000 ALTER TABLE `pools` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `session`
--

DROP TABLE IF EXISTS `session`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `session` (
  `session_id` int NOT NULL AUTO_INCREMENT,
  `machine_id` int NOT NULL DEFAULT '0',
  `pool_id` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`session_id`),
  KEY `fk_machine_id` (`machine_id`),
  KEY `fk_pool_id` (`pool_id`),
  CONSTRAINT `fk_machine_id` FOREIGN KEY (`machine_id`) REFERENCES `machine` (`machine_id`),
  CONSTRAINT `fk_pool_id` FOREIGN KEY (`pool_id`) REFERENCES `pools` (`pool_id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `session`
--

LOCK TABLES `session` WRITE;
/*!40000 ALTER TABLE `session` DISABLE KEYS */;
INSERT INTO `session` VALUES (1,1,1),(2,1,28),(3,1,29),(4,1,30),(5,1,31),(6,1,32),(7,1,33),(8,1,34),(9,1,35),(10,1,36);
/*!40000 ALTER TABLE `session` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_role`
--

DROP TABLE IF EXISTS `user_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_role` (
  `role_id` int NOT NULL AUTO_INCREMENT,
  `role_name` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_role`
--

LOCK TABLES `user_role` WRITE;
/*!40000 ALTER TABLE `user_role` DISABLE KEYS */;
INSERT INTO `user_role` VALUES (1,'admin'),(2,'basic');
/*!40000 ALTER TABLE `user_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `user_id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `hashed_password` char(60) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created` datetime NOT NULL,
  `role_id` int NOT NULL DEFAULT '2',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `users_uc_email` (`email`),
  KEY `role_id` (`role_id`),
  CONSTRAINT `users_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `user_role` (`role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'name1','email1','12345678901234567890123456789012345678901234567890','2022-05-15 17:46:00',2),(2,'mihai22','miahi.mtb10@gamil.com','$2a$12$h6fSr/raQbxO2nXEaJaKzOUJ6yPj0Nae9cWendO8iHIpWtqFRWlIW','2022-05-15 15:30:17',1),(3,'mihai','email1@gmail.com','$2a$12$cmOKKuDUbBrIFc/skt7L0eS.0RYJuxM1sv.Ls.YF.vcWzOcQaf366','2022-05-20 08:57:46',1);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vote`
--

DROP TABLE IF EXISTS `vote`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vote` (
  `vote_id` int NOT NULL AUTO_INCREMENT,
  `pool_id` int NOT NULL DEFAULT '0',
  `option_id` int NOT NULL DEFAULT '0',
  `machine_id` int NOT NULL DEFAULT '0',
  `phone` varchar(15) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`vote_id`),
  KEY `pool_id` (`pool_id`),
  KEY `option_id` (`option_id`),
  KEY `machine_id` (`machine_id`),
  CONSTRAINT `vote_ibfk_1` FOREIGN KEY (`pool_id`) REFERENCES `pools` (`pool_id`),
  CONSTRAINT `vote_ibfk_2` FOREIGN KEY (`option_id`) REFERENCES `pool_option` (`option_id`),
  CONSTRAINT `vote_ibfk_3` FOREIGN KEY (`machine_id`) REFERENCES `machine` (`machine_id`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vote`
--

LOCK TABLES `vote` WRITE;
/*!40000 ALTER TABLE `vote` DISABLE KEYS */;
INSERT INTO `vote` VALUES (1,33,19,1,'0744265634'),(2,35,21,1,'0744265634'),(3,35,21,1,'0744265634'),(4,35,21,1,'0744265634'),(5,35,23,1,'0744265634'),(6,35,23,1,'0744265634'),(7,35,23,1,'0744265634'),(8,36,26,1,'0744265634'),(9,36,25,1,'0744265634'),(10,36,25,1,'0744265634'),(11,36,25,1,'0744265634');
/*!40000 ALTER TABLE `vote` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-05-30 20:43:42