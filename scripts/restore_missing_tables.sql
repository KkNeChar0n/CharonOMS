mysqldump: [Warning] Using a password on the command line interface can be insecure.

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
DROP TABLE IF EXISTS `activity`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `activity` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '?ID',
  `name` varchar(100) NOT NULL COMMENT '??',
  `template_id` int NOT NULL COMMENT '????ģ??ID',
  `start_time` datetime NOT NULL COMMENT '??ʼʱ?',
  `end_time` datetime NOT NULL COMMENT '????ʱ?',
  `status` int DEFAULT '0' COMMENT '״̬??0???? 1?',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '????ʱ?',
  PRIMARY KEY (`id`),
  KEY `template_id` (`template_id`),
  KEY `idx_time_status` (`start_time`,`end_time`,`status`),
  CONSTRAINT `activity_ibfk_1` FOREIGN KEY (`template_id`) REFERENCES `activity_template` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='??';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `activity_detail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `activity_detail` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `activity_id` int NOT NULL COMMENT '?ID',
  `threshold_amount` decimal(10,2) NOT NULL COMMENT '?????????XX??',
  `discount_value` decimal(10,2) NOT NULL COMMENT '?Ż?ֵ????/??/??????ֵ??',
  PRIMARY KEY (`id`),
  KEY `activity_id` (`activity_id`),
  CONSTRAINT `activity_detail_ibfk_1` FOREIGN KEY (`activity_id`) REFERENCES `activity` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='???ϸ???????۹??';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `activity_template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `activity_template` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ģ??ID',
  `name` varchar(100) NOT NULL COMMENT 'ģ???',
  `type` int NOT NULL COMMENT '???ͣ?1???? 2???? 3?',
  `select_type` int NOT NULL COMMENT 'ѡ????ʽ??1?????? 2????Ʒ',
  `status` int DEFAULT '0' COMMENT '״̬??0???? 1?',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '????ʱ?',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '????ʱ?',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='?ģ???';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `activity_template_goods`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `activity_template_goods` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `template_id` int NOT NULL COMMENT 'ģ??ID',
  `goods_id` int DEFAULT NULL COMMENT '??ƷID??select_type=2ʱʹ?ã?',
  `classify_id` int DEFAULT NULL COMMENT '????ID??select_type=1ʱʹ?ã?',
  PRIMARY KEY (`id`),
  KEY `template_id` (`template_id`),
  KEY `goods_id` (`goods_id`),
  KEY `classify_id` (`classify_id`),
  CONSTRAINT `activity_template_goods_ibfk_1` FOREIGN KEY (`template_id`) REFERENCES `activity_template` (`id`) ON DELETE CASCADE,
  CONSTRAINT `activity_template_goods_ibfk_2` FOREIGN KEY (`goods_id`) REFERENCES `goods` (`id`) ON DELETE CASCADE,
  CONSTRAINT `activity_template_goods_ibfk_3` FOREIGN KEY (`classify_id`) REFERENCES `classify` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='?ģ????Ʒ?????';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `childorders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `childorders` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '瀛愯?鍗旾D',
  `parentsid` int NOT NULL COMMENT '鐖惰?鍗旾D',
  `goodsid` int NOT NULL COMMENT '鍟嗗搧ID',
  `amount_receivable` decimal(10,2) NOT NULL COMMENT '搴旀敹閲戦?锛堝晢鍝佹?浠凤級',
  `amount_received` decimal(10,2) NOT NULL COMMENT '瀹炴敹閲戦?锛堟爣鍑嗗敭浠凤級',
  `discount_amount` decimal(10,2) DEFAULT '0.00' COMMENT '?Żݽ',
  `status` int DEFAULT '10' COMMENT '鐘舵?锛?0鑽夌? 20瀹℃牳涓?30宸查?杩?40宸查┏鍥?99宸蹭綔搴',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '鍒涘缓鏃堕棿',
  PRIMARY KEY (`id`),
  KEY `parentsid` (`parentsid`),
  KEY `goodsid` (`goodsid`),
  CONSTRAINT `childorders_ibfk_1` FOREIGN KEY (`parentsid`) REFERENCES `orders` (`id`) ON DELETE CASCADE,
  CONSTRAINT `childorders_ibfk_2` FOREIGN KEY (`goodsid`) REFERENCES `goods` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=72 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='瀛愪骇鍝佽?鍗曡〃';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `orders` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '璁㈠崟ID',
  `student_id` int NOT NULL COMMENT '瀛︾敓ID',
  `expected_payment_time` datetime DEFAULT NULL COMMENT 'Ԥ?Ƹ???ʱ?',
  `amount_receivable` decimal(10,2) NOT NULL COMMENT '搴旀敹閲戦?',
  `amount_received` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '瀹炴敹閲戯拷锟斤拷',
  `discount_amount` decimal(10,2) DEFAULT '0.00' COMMENT '?Żݽ',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '鍒涘缓鏃堕棿',
  `status` tinyint NOT NULL DEFAULT '10' COMMENT '璁㈠崟鐘舵?: 10=鑽夌?, 20=瀹℃牳涓? 30=宸查?杩? 40=宸查┏鍥? 99=宸蹭綔搴',
  PRIMARY KEY (`id`),
  KEY `idx_student_id` (`student_id`),
  KEY `idx_status` (`status`),
  KEY `idx_create_time` (`create_time`),
  CONSTRAINT `fk_order_student` FOREIGN KEY (`student_id`) REFERENCES `student` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='璁㈠崟琛';
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `orders_activity`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `orders_activity` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `orders_id` int NOT NULL COMMENT '????ID',
  `activity_id` int NOT NULL COMMENT '?ID',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '????ʱ?',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_orders_activity` (`orders_id`,`activity_id`),
  KEY `idx_orders_id` (`orders_id`),
  KEY `idx_activity_id` (`activity_id`),
  CONSTRAINT `orders_activity_ibfk_1` FOREIGN KEY (`orders_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE,
  CONSTRAINT `orders_activity_ibfk_2` FOREIGN KEY (`activity_id`) REFERENCES `activity` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='??????????';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

