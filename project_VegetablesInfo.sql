-- MySQL dump 10.13  Distrib 5.7.31, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: project
-- ------------------------------------------------------
-- Server version	5.7.31-0ubuntu0.16.04.1

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
-- Table structure for table `VegetablesInfo`
--

DROP TABLE IF EXISTS `VegetablesInfo`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `VegetablesInfo` (
  `Veid` int(11) NOT NULL,
  `Vename` varchar(30) NOT NULL,
  `Vedata` int(11) NOT NULL,
  PRIMARY KEY (`Veid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `VegetablesInfo`
--

LOCK TABLES `VegetablesInfo` WRITE;
/*!40000 ALTER TABLE `VegetablesInfo` DISABLE KEYS */;
INSERT INTO `VegetablesInfo` VALUES (1,'Potato',6),(2,'Tomato',13),(3,'Carrot',28),(4,'Onion',36),(5,'Corn',45),(6,'Cucumber',49),(7,'Pepper',9),(8,'Radish',17),(9,'Potato',50),(10,'Tomato',51),(11,'Carrot',62),(12,'Onion',74),(13,'Corn',85),(14,'Cucumber',88),(15,'Pepper',90),(16,'Radish',99),(17,'Potato',100),(18,'Tomato',169),(19,'Carrot',146),(20,'Onion',122),(21,'Corn',183),(22,'Cucumber',201),(23,'Pepper',199),(24,'Radish',287),(25,'Potato',243),(26,'Tomato',291),(27,'Carrot',299),(28,'Onion',300),(29,'Corn',357),(30,'Cucumber',288),(31,'Pepper',477),(32,'Radish',499),(33,'Potato',100),(34,'Tomato',169),(35,'Carrot',146),(36,'Onion',122),(37,'Corn',183),(38,'Cucumber',201),(39,'Pepper',199),(40,'Radish',287),(41,'Potato',243),(42,'Tomato',291),(43,'Carrot',299),(44,'Onion',300),(45,'Corn',357),(46,'Cucumber',288),(47,'Pepper',477),(48,'Radish',499),(49,'Potato',100),(50,'Tomato',169),(51,'Carrot',146),(52,'Onion',122),(53,'Corn',183),(54,'Cucumber',201),(55,'Pepper',199),(56,'Radish',287),(57,'Potato',243),(58,'Tomato',291),(59,'Carrot',299),(60,'Onion',300),(61,'Corn',357),(62,'Cucumber',288),(63,'Pepper',477),(64,'Radish',499);
/*!40000 ALTER TABLE `VegetablesInfo` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-07-29 10:14:10
