/*M!999999\- enable the sandbox mode */
-- MariaDB dump 10.19  Distrib 10.11.14-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: newacct
-- ------------------------------------------------------
-- Server version       10.11.14-MariaDB-0ubuntu0.24.04.1

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
-- Table structure for table `category`
--

DROP TABLE IF EXISTS `category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `category` (
                            `category_cd` varchar(64) NOT NULL,
                            PRIMARY KEY (`category_cd`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `category_rules`
--

DROP TABLE IF EXISTS `category_rules`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `category_rules` (
                                  `id` varchar(64) NOT NULL,
                                  `user_id` varchar(64) NOT NULL,
                                  `pattern` varchar(128) NOT NULL,
                                  `category_cd` varchar(128) NOT NULL,
                                  `comments` varchar(255) DEFAULT NULL,
                                  PRIMARY KEY (`id`),
                                  KEY `fk_rules_category` (`category_cd`),
                                  KEY `fk_rules_user` (`user_id`),
                                  CONSTRAINT `fk_rules_category` FOREIGN KEY (`category_cd`) REFERENCES `category` (`category_cd`) ON DELETE NO ACTION ON UPDATE NO ACTION,
                                  CONSTRAINT `fk_rules_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `transaction`
--

DROP TABLE IF EXISTS `transaction`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `transaction` (
                               `id` varchar(64) NOT NULL,
                               `user_id` varchar(64) NOT NULL,
                               `bank_id` varchar(128) DEFAULT NULL,
                               `account_id` varchar(128) DEFAULT NULL,
                               `account_type` varchar(128) DEFAULT NULL,
                               `transaction_type` varchar(128) DEFAULT NULL,
                               `date_posted` datetime DEFAULT NULL,
                               `amount` decimal(14,2) NOT NULL,
                               `bank_transaction_id` varchar(128) DEFAULT NULL,
                               `name` varchar(255) DEFAULT NULL,
                               `category_cd` varchar(128) DEFAULT NULL,
                               `comments` varchar(255) DEFAULT NULL,
                               `hash` varchar(255) DEFAULT NULL,
                               PRIMARY KEY (`id`),
                               KEY `fk_categ_trans` (`category_cd`),
                               KEY `fk_user_trans` (`user_id`),
                               CONSTRAINT `fk_categ_trans` FOREIGN KEY (`category_cd`) REFERENCES `category` (`category_cd`) ON DELETE NO ACTION ON UPDATE NO ACTION,
                               CONSTRAINT `fk_user_trans` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
                        `user_id` varchar(64) NOT NULL,
                        `name` varchar(255) NOT NULL,
                        `login` varchar(128) NOT NULL,
                        `password` varchar(64) DEFAULT NULL,
                        `connect_ind` varchar(1) DEFAULT NULL,
                        `admin_ind` varchar(1) DEFAULT NULL,
                        PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-11-14 17:32:27
